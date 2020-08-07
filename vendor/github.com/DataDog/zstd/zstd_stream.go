package zstd

/*
#define ZSTD_STATIC_LINKING_ONLY
#include "stdint.h"  // for uintptr_t
#include "zstd.h"

typedef struct compressStream2_result_s {
	size_t return_code;
	size_t bytes_consumed;
	size_t bytes_written;
} compressStream2_result;

static void ZSTD_compressStream2_wrapper(compressStream2_result* result, ZSTD_CCtx* ctx, uintptr_t dst, size_t maxDstSize, const uintptr_t src, size_t srcSize) {
	ZSTD_outBuffer outBuffer = { (void*)dst, maxDstSize, 0 };
	ZSTD_inBuffer inBuffer = { (void*)src, srcSize, 0 };
	size_t retCode = ZSTD_compressStream2(ctx, &outBuffer, &inBuffer, ZSTD_e_continue);

	result->return_code = retCode;
	result->bytes_consumed = inBuffer.pos;
	result->bytes_written = outBuffer.pos;
}

static void ZSTD_compressStream2_finish(compressStream2_result* result, ZSTD_CCtx* ctx, uintptr_t dst, size_t maxDstSize, const uintptr_t src, size_t srcSize) {
	ZSTD_outBuffer outBuffer = { (void*)dst, maxDstSize, 0 };
	ZSTD_inBuffer inBuffer = { (void*)src, srcSize, 0 };
	size_t retCode = ZSTD_compressStream2(ctx, &outBuffer, &inBuffer, ZSTD_e_end);

	result->return_code = retCode;
	result->bytes_consumed = inBuffer.pos;
	result->bytes_written = outBuffer.pos;
}

// decompressStream2_result is the same as compressStream2_result, but keep 2 separate struct for easier changes
typedef struct decompressStream2_result_s {
	size_t return_code;
	size_t bytes_consumed;
	size_t bytes_written;
} decompressStream2_result;

static void ZSTD_decompressStream_wrapper(decompressStream2_result* result, ZSTD_DCtx* ctx, uintptr_t dst, size_t maxDstSize, const uintptr_t src, size_t srcSize) {
	ZSTD_outBuffer outBuffer = { (void*)dst, maxDstSize, 0 };
	ZSTD_inBuffer inBuffer = { (void*)src, srcSize, 0 };
	size_t retCode = ZSTD_decompressStream(ctx, &outBuffer, &inBuffer);

	result->return_code = retCode;
	result->bytes_consumed = inBuffer.pos;
	result->bytes_written = outBuffer.pos;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"sync"
	"unsafe"
)

var errShortRead = errors.New("short read")
var errReaderClosed = errors.New("Reader is closed")

// Writer is an io.WriteCloser that zstd-compresses its input.
type Writer struct {
	CompressionLevel int

	ctx              *C.ZSTD_CCtx
	dict             []byte
	srcBuffer        []byte
	dstBuffer        []byte
	firstError       error
	underlyingWriter io.Writer
	resultBuffer     *C.compressStream2_result
}

func resize(in []byte, newSize int) []byte {
	if in == nil {
		return make([]byte, newSize)
	}
	if newSize <= cap(in) {
		return in[:newSize]
	}
	toAdd := newSize - len(in)
	return append(in, make([]byte, toAdd)...)
}

// NewWriter creates a new Writer with default compression options.  Writes to
// the writer will be written in compressed form to w.
func NewWriter(w io.Writer) *Writer {
	return NewWriterLevelDict(w, DefaultCompression, nil)
}

// NewWriterLevel is like NewWriter but specifies the compression level instead
// of assuming default compression.
//
// The level can be DefaultCompression or any integer value between BestSpeed
// and BestCompression inclusive.
func NewWriterLevel(w io.Writer, level int) *Writer {
	return NewWriterLevelDict(w, level, nil)

}

// NewWriterLevelDict is like NewWriterLevel but specifies a dictionary to
// compress with.  If the dictionary is empty or nil it is ignored. The dictionary
// should not be modified until the writer is closed.
func NewWriterLevelDict(w io.Writer, level int, dict []byte) *Writer {
	var err error
	ctx := C.ZSTD_createCStream()

	// Load dictionnary if any
	if dict != nil {
		err = getError(int(C.ZSTD_CCtx_loadDictionary(ctx,
			unsafe.Pointer(&dict[0]),
			C.size_t(len(dict)),
		)))
	}

	if err == nil {
		// Only set level if the ctx is not in error already
		err = getError(int(C.ZSTD_CCtx_setParameter(ctx, C.ZSTD_c_compressionLevel, C.int(level))))
	}

	return &Writer{
		CompressionLevel: level,
		ctx:              ctx,
		dict:             dict,
		srcBuffer:        make([]byte, 0),
		dstBuffer:        make([]byte, CompressBound(1024)),
		firstError:       err,
		underlyingWriter: w,
		resultBuffer:     new(C.compressStream2_result),
	}
}

// Write writes a compressed form of p to the underlying io.Writer.
func (w *Writer) Write(p []byte) (int, error) {
	if w.firstError != nil {
		return 0, w.firstError
	}
	if len(p) == 0 {
		return 0, nil
	}
	// Check if dstBuffer is enough
	w.dstBuffer = w.dstBuffer[0:cap(w.dstBuffer)]
	if len(w.dstBuffer) < CompressBound(len(p)) {
		w.dstBuffer = make([]byte, CompressBound(len(p)))
	}

	// Do not do an extra memcopy if zstd ingest all input data
	srcData := p
	fastPath := len(w.srcBuffer) == 0
	if !fastPath {
		w.srcBuffer = append(w.srcBuffer, p...)
		srcData = w.srcBuffer
	}

	srcPtr := C.uintptr_t(uintptr(0)) // Do not point anywhere, if src is empty
	if len(srcData) > 0 {
		srcPtr = C.uintptr_t(uintptr(unsafe.Pointer(&srcData[0])))
	}

	C.ZSTD_compressStream2_wrapper(
		w.resultBuffer,
		w.ctx,
		C.uintptr_t(uintptr(unsafe.Pointer(&w.dstBuffer[0]))),
		C.size_t(len(w.dstBuffer)),
		srcPtr,
		C.size_t(len(srcData)),
	)
	runtime.KeepAlive(p) // Ensure p is kept until here so pointer doesn't disappear during C call
	ret := int(w.resultBuffer.return_code)
	if err := getError(ret); err != nil {
		return 0, err
	}

	consumed := int(w.resultBuffer.bytes_consumed)
	if !fastPath {
		w.srcBuffer = w.srcBuffer[consumed:]
	} else {
		remaining := len(p) - consumed
		if remaining > 0 {
			// We still have some non-consumed data, copy remaining data to srcBuffer
			// Try to not reallocate w.srcBuffer if we already have enough space
			if cap(w.srcBuffer) >= remaining {
				w.srcBuffer = w.srcBuffer[0:remaining]
			} else {
				w.srcBuffer = make([]byte, remaining)
			}
			copy(w.srcBuffer, p[consumed:])
		}
	}

	written := int(w.resultBuffer.bytes_written)
	// Write to underlying buffer
	_, err := w.underlyingWriter.Write(w.dstBuffer[:written])

	// Same behaviour as zlib, we can't know how much data we wrote, only
	// if there was an error
	if err != nil {
		return 0, err
	}
	return len(p), err
}

// Close closes the Writer, flushing any unwritten data to the underlying
// io.Writer and freeing objects, but does not close the underlying io.Writer.
func (w *Writer) Close() error {
	if w.firstError != nil {
		return w.firstError
	}

	ret := 1 // So we loop at least once
	for ret > 0 {
		srcPtr := C.uintptr_t(uintptr(0)) // Do not point anywhere, if src is empty
		if len(w.srcBuffer) > 0 {
			srcPtr = C.uintptr_t(uintptr(unsafe.Pointer(&w.srcBuffer[0])))
		}

		C.ZSTD_compressStream2_finish(
			w.resultBuffer,
			w.ctx,
			C.uintptr_t(uintptr(unsafe.Pointer(&w.dstBuffer[0]))),
			C.size_t(len(w.dstBuffer)),
			srcPtr,
			C.size_t(len(w.srcBuffer)),
		)
		ret = int(w.resultBuffer.return_code)
		if err := getError(ret); err != nil {
			return err
		}
		w.srcBuffer = w.srcBuffer[w.resultBuffer.bytes_consumed:]
		written := int(w.resultBuffer.bytes_written)
		w.underlyingWriter.Write(w.dstBuffer[:written])

		if ret > 0 { // We have a hint if we need to resize the dstBuffer
			w.dstBuffer = w.dstBuffer[:cap(w.dstBuffer)]
			if len(w.dstBuffer) < ret {
				w.dstBuffer = make([]byte, ret)
			}
		}
	}

	return getError(int(C.ZSTD_freeCStream(w.ctx)))
}

// cSize is the recommended size of reader.compressionBuffer. This func and
// invocation allow for a one-time check for validity.
var cSize = func() int {
	v := int(C.ZSTD_DStreamInSize())
	if v <= 0 {
		panic(fmt.Errorf("ZSTD_DStreamInSize() returned invalid size: %v", v))
	}
	return v
}()

// dSize is the recommended size of reader.decompressionBuffer. This func and
// invocation allow for a one-time check for validity.
var dSize = func() int {
	v := int(C.ZSTD_DStreamOutSize())
	if v <= 0 {
		panic(fmt.Errorf("ZSTD_DStreamOutSize() returned invalid size: %v", v))
	}
	return v
}()

// cPool is a pool of buffers for use in reader.compressionBuffer. Buffers are
// taken from the pool in NewReaderDict, returned in reader.Close(). Returns a
// pointer to a slice to avoid the extra allocation of returning the slice as a
// value.
var cPool = sync.Pool{
	New: func() interface{} {
		buff := make([]byte, cSize)
		return &buff
	},
}

// dPool is a pool of buffers for use in reader.decompressionBuffer. Buffers are
// taken from the pool in NewReaderDict, returned in reader.Close(). Returns a
// pointer to a slice to avoid the extra allocation of returning the slice as a
// value.
var dPool = sync.Pool{
	New: func() interface{} {
		buff := make([]byte, dSize)
		return &buff
	},
}

// reader is an io.ReadCloser that decompresses when read from.
type reader struct {
	ctx                 *C.ZSTD_DCtx
	compressionBuffer   []byte
	compressionLeft     int
	decompressionBuffer []byte
	decompOff           int
	decompSize          int
	dict                []byte
	firstError          error
	recommendedSrcSize  int
	resultBuffer        *C.decompressStream2_result
	underlyingReader    io.Reader
}

// NewReader creates a new io.ReadCloser.  Reads from the returned ReadCloser
// read and decompress data from r.  It is the caller's responsibility to call
// Close on the ReadCloser when done.  If this is not done, underlying objects
// in the zstd library will not be freed.
func NewReader(r io.Reader) io.ReadCloser {
	return NewReaderDict(r, nil)
}

// NewReaderDict is like NewReader but uses a preset dictionary.  NewReaderDict
// ignores the dictionary if it is nil.
func NewReaderDict(r io.Reader, dict []byte) io.ReadCloser {
	var err error
	ctx := C.ZSTD_createDStream()
	if len(dict) == 0 {
		err = getError(int(C.ZSTD_initDStream(ctx)))
	} else {
		err = getError(int(C.ZSTD_DCtx_reset(ctx, C.ZSTD_reset_session_only)))
		if err == nil {
			// Only load dictionary if we succesfully inited the context
			err = getError(int(C.ZSTD_DCtx_loadDictionary(
				ctx,
				unsafe.Pointer(&dict[0]),
				C.size_t(len(dict)))))
		}
	}
	compressionBufferP := cPool.Get().(*[]byte)
	decompressionBufferP := dPool.Get().(*[]byte)
	return &reader{
		ctx:                 ctx,
		dict:                dict,
		compressionBuffer:   *compressionBufferP,
		decompressionBuffer: *decompressionBufferP,
		firstError:          err,
		recommendedSrcSize:  cSize,
		resultBuffer:        new(C.decompressStream2_result),
		underlyingReader:    r,
	}
}

// Close frees the allocated C objects
func (r *reader) Close() error {
	if r.firstError != nil {
		return r.firstError
	}

	cb := r.compressionBuffer
	db := r.decompressionBuffer
	// Ensure that we won't resuse buffer
	r.firstError = errReaderClosed
	r.compressionBuffer = nil
	r.decompressionBuffer = nil

	cPool.Put(&cb)
	dPool.Put(&db)
	return getError(int(C.ZSTD_freeDStream(r.ctx)))
}

func (r *reader) Read(p []byte) (int, error) {
	if r.firstError != nil {
		return 0, r.firstError
	}

	// If we already have enough bytes, return
	if r.decompSize-r.decompOff >= len(p) {
		copy(p, r.decompressionBuffer[r.decompOff:])
		r.decompOff += len(p)
		return len(p), nil
	}

	copy(p, r.decompressionBuffer[r.decompOff:r.decompSize])
	got := r.decompSize - r.decompOff
	r.decompSize = 0
	r.decompOff = 0

	for got < len(p) {
		// Populate src
		src := r.compressionBuffer
		reader := r.underlyingReader
		n, err := TryReadFull(reader, src[r.compressionLeft:])
		if err != nil && err != errShortRead { // Handle underlying reader errors first
			return 0, fmt.Errorf("failed to read from underlying reader: %s", err)
		} else if n == 0 && r.compressionLeft == 0 {
			return got, io.EOF
		}
		src = src[:r.compressionLeft+n]

		// C code
		srcPtr := C.uintptr_t(uintptr(0)) // Do not point anywhere, if src is empty
		if len(src) > 0 {
			srcPtr = C.uintptr_t(uintptr(unsafe.Pointer(&src[0])))
		}

		C.ZSTD_decompressStream_wrapper(
			r.resultBuffer,
			r.ctx,
			C.uintptr_t(uintptr(unsafe.Pointer(&r.decompressionBuffer[0]))),
			C.size_t(len(r.decompressionBuffer)),
			srcPtr,
			C.size_t(len(src)),
		)
		retCode := int(r.resultBuffer.return_code)

		// Keep src here eventhough we reuse later, the code might be deleted at some point
		runtime.KeepAlive(src)
		if err = getError(retCode); err != nil {
			return 0, fmt.Errorf("failed to decompress: %s", err)
		}

		// Put everything in buffer
		bytesConsumed := int(r.resultBuffer.bytes_consumed)
		if bytesConsumed < len(src) {
			left := src[bytesConsumed:]
			copy(r.compressionBuffer, left)
		}
		r.compressionLeft = len(src) - int(bytesConsumed)
		r.decompSize = int(r.resultBuffer.bytes_written)
		r.decompOff = copy(p[got:], r.decompressionBuffer[:r.decompSize])
		got += r.decompOff

		// Resize buffers
		nsize := retCode // Hint for next src buffer size
		if nsize <= 0 {
			// Reset to recommended size
			nsize = r.recommendedSrcSize
		}
		if nsize < r.compressionLeft {
			nsize = r.compressionLeft
		}
		r.compressionBuffer = resize(r.compressionBuffer, nsize)
	}
	return got, nil
}

// TryReadFull reads buffer just as ReadFull does
// Here we expect that buffer may end and we do not return ErrUnexpectedEOF as ReadAtLeast does.
// We return errShortRead instead to distinguish short reads and failures.
// We cannot use ReadFull/ReadAtLeast because it masks Reader errors, such as network failures
// and causes panic instead of error.
func TryReadFull(r io.Reader, buf []byte) (n int, err error) {
	for n < len(buf) && err == nil {
		var nn int
		nn, err = r.Read(buf[n:])
		n += nn
	}
	if n == len(buf) && err == io.EOF {
		err = nil // EOF at the end is somewhat expected
	} else if err == io.EOF {
		err = errShortRead
	}
	return
}
