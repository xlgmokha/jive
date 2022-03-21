#!/bin/bash

__jive_root_dir="$(builtin cd "$(\dirname "${BASH_SOURCE[0]}")" && \pwd)"
__jive_script="${__jive_root_dir}/jive.sh"

__mtime_of_jive_script="$(\date -r "${__jive_script}" +%s)"
__jive_auto_reload() {
  local current_mtime
  current_mtime="$(\date -r "${__jive_script}" +%s)"

  if [[ "${current_mtime}" != "${__mtime_of_jive_script}" ]]; then
    echo "Reloading... ${__jive_script}"
    . "${__jive_script}"
  fi
}

__jive_exec() {
  cd "$HOME/src/github.com/xlgmokha/jive/" || exit 1
  go run main.go "$@"
}

__jive_open_pipe() {
  local tmpfile
  tmpfile="$(\mktemp -u)"

  exec 42>"${tmpfile}" # Open the tempfile for writing on FD 42.
  exec 8<"${tmpfile}" # Open the tempfile for reading on FD 8.
  \rm -f "${tmpfile}" # Unlink the tempfile. (we've already opened it).
}

__jive_execute_task() {
  local task=$1

  case "${task}" in
    cd:*)
      # shellcheck disable=SC2164
      cd "${task//cd:/}"
      ;;
    ctags:*)
      # shellcheck disable=SC2164
      ctags -R "${task//ctags:/}"
      ;;
    setenv:*)
      export "${task//setenv:/}"
      ;;
    *)
      echo "Woof! ${task}"
      ;;
  esac
}

__jive_flush_tasks() {
  local task
  while \read -r task; do
    __jive_execute_task "${task}"
  done <&8

  __jive_close_pipe
}

__jive_close_pipe() {
  exec 8<&- # close FD 8.
  exec 42<&- # close FD 42.
}

jive() {
  __jive_auto_reload
  __jive_open_pipe
  __jive_exec "$@"
  __jive_flush_tasks
}
