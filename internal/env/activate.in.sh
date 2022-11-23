# This script should be sourced.

# In this script, we should assume a POSIX shell. Special options can be set, so we should be careful about
# unset variables and commands returning non-0 status.

if [ -n "${GOENV_PATH:-}" ]; then
  echo "An environment is already active (deactivate with \"deactivate\")." >&2
else
  deactivate() {
    PATH="$GOENV_SAVED_PATH"

    GOROOT="$GOENV_SAVED_GOROOT"
    if [ -z "$GOROOT" ]; then
      unset GOROOT
    fi

    GOPATH="$GOENV_SAVED_GOPATH"
    if [ -z "$GOPATH" ]; then
      unset GOPATH
    fi

    unset GOENV_PATH
    unset GOENV_SAVED_PATH
    unset GOENV_SAVED_GOROOT
    unset GOENV_SAVED_GOPATH
    unset -f deactivate
  }

  export GOENV_PATH={{ .Dir.Escape }}

  GOENV_SAVED_GOROOT="${GOROOT:-}"
  export GOROOT={{ .GoRoot.Escape }}
  echo GOROOT={{ .GoRoot.Escape.Escape }}

  GOENV_SAVED_GOPATH="${GOPATH:-}"
  export GOPATH={{ .GoPath.Escape }}
  echo GOPATH={{ .GoPath.Escape.Escape }}

  GOENV_SAVED_PATH="$PATH"
  PATH="$GOPATH/bin:$GOROOT/bin:$PATH"

  echo "Deactivate with: deactivate"
fi
