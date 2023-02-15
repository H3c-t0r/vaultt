package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/eventbus"
	"google.golang.org/protobuf/encoding/protojson"
	"nhooyr.io/websocket"
)

type eventSubscribeArgs struct {
	ctx       context.Context
	logger    hclog.Logger
	events    *eventbus.EventBus
	ns        *namespace.Namespace
	eventType logical.EventType
	conn      *websocket.Conn
	json      bool
}

// handleEventsSubscribeWebsocket runs forever, returning a websocket error code and reason
// only if the connection closes or there was an error.
func handleEventsSubscribeWebsocket(args eventSubscribeArgs) (websocket.StatusCode, string, error) {
	ctx := args.ctx
	logger := args.logger
	ch, cancel, err := args.events.Subscribe(ctx, args.ns, args.eventType)
	if err != nil {
		logger.Info("Error subscribing", "error", err)
		return websocket.StatusUnsupportedData, "Error subscribing", nil
	}
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Websocket context is done, closing the connection")
			return websocket.StatusNormalClosure, "", nil
		case message := <-ch:
			logger.Debug("Sending message to websocket", "message", message)
			var messageBytes []byte
			if args.json {
				messageBytes, err = protojson.Marshal(message)
			} else {
				messageBytes, err = proto.Marshal(message)
			}
			if err != nil {
				logger.Warn("Could not serialize websocket event", "error", err)
				return 0, "", err
			}
			messageString := string(messageBytes) + "\n"
			err = args.conn.Write(ctx, websocket.MessageText, []byte(messageString))
			if err != nil {
				return 0, "", err
			}
		}
	}
}

func handleEventsSubscribe(core *vault.Core, req *logical.Request) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := core.Logger().Named("events-subscribe")
		logger.Debug("Got request to", "url", r.URL, "version", r.Proto)

		ctx := r.Context()

		// ACL check
		_, _, err := core.CheckToken(ctx, req, false)
		if err != nil {
			if errors.Is(err, logical.ErrPermissionDenied) {
				respondError(w, http.StatusUnauthorized, logical.ErrPermissionDenied)
				return
			}
			logger.Debug("Error validating token", "error", err)
			respondError(w, http.StatusInternalServerError, fmt.Errorf("error validating token"))
			return
		}

		ns, err := namespace.FromContext(ctx)
		if err != nil {
			logger.Info("Could not find namespace", "error", err)
			respondError(w, http.StatusInternalServerError, fmt.Errorf("could not find namespace"))
			return
		}

		prefix := "/v1/sys/events/subscribe/"
		if ns.ID != namespace.RootNamespaceID {
			prefix = fmt.Sprintf("/v1/%ssys/events/subscribe/", ns.Path)
		}
		eventTypeStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
		if eventTypeStr == "" {
			respondError(w, http.StatusBadRequest, fmt.Errorf("did not specify eventType to subscribe to"))
			return
		}
		eventType := logical.EventType(eventTypeStr)

		json := false
		jsonRaw := r.URL.Query().Get("json")
		if jsonRaw != "" {
			var err error
			json, err = strconv.ParseBool(jsonRaw)
			if err != nil {
				respondError(w, http.StatusBadRequest, fmt.Errorf("invalid parameter for JSON: %v", jsonRaw))
				return
			}
		}

		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			logger.Info("Could not accept as websocket", "error", err)
			respondError(w, http.StatusInternalServerError, fmt.Errorf("could not accept as websocket"))
			return
		}

		// we don't expect any incoming messages
		ctx = conn.CloseRead(ctx)
		// start the pinger
		go func() {
			for {
				time.Sleep(30 * time.Second) // not too aggressive, but keep the HTTP connection alive
				err := conn.Ping(ctx)
				if err != nil {
					return
				}
			}
		}()

		closeStatus, closeReason, err := handleEventsSubscribeWebsocket(eventSubscribeArgs{ctx, logger, core.Events(), ns, eventType, conn, json})
		if err != nil {
			closeStatus = websocket.CloseStatus(err)
			if closeStatus == -1 {
				closeStatus = websocket.StatusInternalError
			}
			closeReason = fmt.Sprintf("Internal error: %v", err)
			logger.Debug("Error from websocket handler", "error", err)
		}
		// Close() will panic if the reason is greater than this length
		if len(closeReason) > 123 {
			logger.Debug("Truncated close reason", "closeReason", closeReason)
			closeReason = closeReason[:123]
		}
		err = conn.Close(closeStatus, closeReason)
		if err != nil {
			logger.Debug("Error closing websocket", "error", err)
		}
	})
}
