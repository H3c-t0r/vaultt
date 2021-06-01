// Code generated by cmd/cgo -godefs; DO NOT EDIT.
// cgo -godefs sysconf_defs_linux.go

package sysconf

const (
	SC_AIO_LISTIO_MAX               = 0x17
	SC_AIO_MAX                      = 0x18
	SC_AIO_PRIO_DELTA_MAX           = 0x19
	SC_ARG_MAX                      = 0x0
	SC_ATEXIT_MAX                   = 0x57
	SC_BC_BASE_MAX                  = 0x24
	SC_BC_DIM_MAX                   = 0x25
	SC_BC_SCALE_MAX                 = 0x26
	SC_BC_STRING_MAX                = 0x27
	SC_CHILD_MAX                    = 0x1
	SC_CLK_TCK                      = 0x2
	SC_COLL_WEIGHTS_MAX             = 0x28
	SC_DELAYTIMER_MAX               = 0x1a
	SC_EXPR_NEST_MAX                = 0x2a
	SC_GETGR_R_SIZE_MAX             = 0x45
	SC_GETPW_R_SIZE_MAX             = 0x46
	SC_HOST_NAME_MAX                = 0xb4
	SC_IOV_MAX                      = 0x3c
	SC_LINE_MAX                     = 0x2b
	SC_LOGIN_NAME_MAX               = 0x47
	SC_MQ_OPEN_MAX                  = 0x1b
	SC_MQ_PRIO_MAX                  = 0x1c
	SC_NGROUPS_MAX                  = 0x3
	SC_OPEN_MAX                     = 0x4
	SC_PAGE_SIZE                    = 0x1e
	SC_PAGESIZE                     = 0x1e
	SC_THREAD_DESTRUCTOR_ITERATIONS = 0x49
	SC_THREAD_KEYS_MAX              = 0x4a
	SC_THREAD_STACK_MIN             = 0x4b
	SC_THREAD_THREADS_MAX           = 0x4c
	SC_RE_DUP_MAX                   = 0x2c
	SC_RTSIG_MAX                    = 0x1f
	SC_SEM_NSEMS_MAX                = 0x20
	SC_SEM_VALUE_MAX                = 0x21
	SC_SIGQUEUE_MAX                 = 0x22
	SC_STREAM_MAX                   = 0x5
	SC_SYMLOOP_MAX                  = 0xad
	SC_TIMER_MAX                    = 0x23
	SC_TTY_NAME_MAX                 = 0x48
	SC_TZNAME_MAX                   = 0x6

	SC_ADVISORY_INFO              = 0x84
	SC_ASYNCHRONOUS_IO            = 0xc
	SC_BARRIERS                   = 0x85
	SC_CLOCK_SELECTION            = 0x89
	SC_CPUTIME                    = 0x8a
	SC_FSYNC                      = 0xf
	SC_IPV6                       = 0xeb
	SC_JOB_CONTROL                = 0x7
	SC_MAPPED_FILES               = 0x10
	SC_MEMLOCK                    = 0x11
	SC_MEMLOCK_RANGE              = 0x12
	SC_MEMORY_PROTECTION          = 0x13
	SC_MESSAGE_PASSING            = 0x14
	SC_MONOTONIC_CLOCK            = 0x95
	SC_PRIORITIZED_IO             = 0xd
	SC_PRIORITY_SCHEDULING        = 0xa
	SC_RAW_SOCKETS                = 0xec
	SC_READER_WRITER_LOCKS        = 0x99
	SC_REALTIME_SIGNALS           = 0x9
	SC_REGEXP                     = 0x9b
	SC_SAVED_IDS                  = 0x8
	SC_SEMAPHORES                 = 0x15
	SC_SHARED_MEMORY_OBJECTS      = 0x16
	SC_SHELL                      = 0x9d
	SC_SPAWN                      = 0x9f
	SC_SPIN_LOCKS                 = 0x9a
	SC_SPORADIC_SERVER            = 0xa0
	SC_SS_REPL_MAX                = 0xf1
	SC_SYNCHRONIZED_IO            = 0xe
	SC_THREAD_ATTR_STACKADDR      = 0x4d
	SC_THREAD_ATTR_STACKSIZE      = 0x4e
	SC_THREAD_CPUTIME             = 0x8b
	SC_THREAD_PRIO_INHERIT        = 0x50
	SC_THREAD_PRIO_PROTECT        = 0x51
	SC_THREAD_PRIORITY_SCHEDULING = 0x4f
	SC_THREAD_PROCESS_SHARED      = 0x52
	SC_THREAD_ROBUST_PRIO_INHERIT = 0xf7
	SC_THREAD_ROBUST_PRIO_PROTECT = 0xf8
	SC_THREAD_SAFE_FUNCTIONS      = 0x44
	SC_THREAD_SPORADIC_SERVER     = 0xa1
	SC_THREADS                    = 0x43
	SC_TIMEOUTS                   = 0xa4
	SC_TIMERS                     = 0xb
	SC_TRACE                      = 0xb5
	SC_TRACE_EVENT_FILTER         = 0xb6
	SC_TRACE_EVENT_NAME_MAX       = 0xf2
	SC_TRACE_INHERIT              = 0xb7
	SC_TRACE_LOG                  = 0xb8
	SC_TRACE_NAME_MAX             = 0xf3
	SC_TRACE_SYS_MAX              = 0xf4
	SC_TRACE_USER_EVENT_MAX       = 0xf5
	SC_TYPED_MEMORY_OBJECTS       = 0xa5
	SC_VERSION                    = 0x1d

	SC_V7_ILP32_OFF32  = 0xed
	SC_V7_ILP32_OFFBIG = 0xee
	SC_V7_LP64_OFF64   = 0xef
	SC_V7_LPBIG_OFFBIG = 0xf0

	SC_V6_ILP32_OFF32  = 0xb0
	SC_V6_ILP32_OFFBIG = 0xb1
	SC_V6_LP64_OFF64   = 0xb2
	SC_V6_LPBIG_OFFBIG = 0xb3

	SC_2_C_BIND         = 0x2f
	SC_2_C_DEV          = 0x30
	SC_2_C_VERSION      = 0x60
	SC_2_CHAR_TERM      = 0x5f
	SC_2_FORT_DEV       = 0x31
	SC_2_FORT_RUN       = 0x32
	SC_2_LOCALEDEF      = 0x34
	SC_2_PBS            = 0xa8
	SC_2_PBS_ACCOUNTING = 0xa9
	SC_2_PBS_CHECKPOINT = 0xaf
	SC_2_PBS_LOCATE     = 0xaa
	SC_2_PBS_MESSAGE    = 0xab
	SC_2_PBS_TRACK      = 0xac
	SC_2_SW_DEV         = 0x33
	SC_2_UPE            = 0x61
	SC_2_VERSION        = 0x2e

	SC_XOPEN_CRYPT            = 0x5c
	SC_XOPEN_ENH_I18N         = 0x5d
	SC_XOPEN_REALTIME         = 0x82
	SC_XOPEN_REALTIME_THREADS = 0x83
	SC_XOPEN_SHM              = 0x5e
	SC_XOPEN_STREAMS          = 0xf6
	SC_XOPEN_UNIX             = 0x5b
	SC_XOPEN_VERSION          = 0x59
	SC_XOPEN_XCU_VERSION      = 0x5a

	SC_PHYS_PAGES       = 0x55
	SC_AVPHYS_PAGES     = 0x56
	SC_NPROCESSORS_CONF = 0x53
	SC_NPROCESSORS_ONLN = 0x54
	SC_UIO_MAXIOV       = 0x3c
)
