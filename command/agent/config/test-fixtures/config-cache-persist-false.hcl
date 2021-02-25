pid_file = "./pidfile"

cache {
    persist "kubernetes" {
        exit_on_err = false
        keep_after_import = false
        path = "/tmp/bolt-file.db"
    }
}

listener "tcp" {
    address = "127.0.0.1:8300"
    tls_disable = true
}
