// +build !windows,!darwin

package zenity

import (
	"os/exec"
	"strconv"

	"github.com/ncruces/zenity/internal/zenutil"
)

func message(kind messageKind, text string, opts options) (bool, error) {
	var args []string
	switch kind {
	case questionKind:
		args = append(args, "--question")
	case infoKind:
		args = append(args, "--info")
	case warningKind:
		args = append(args, "--warning")
	case errorKind:
		args = append(args, "--error")
	}
	if text != "" {
		args = append(args, "--text", text, "--no-markup")
	}
	if opts.title != nil {
		args = append(args, "--title", *opts.title)
	}
	if opts.width > 0 {
		args = append(args, "--width", strconv.FormatUint(uint64(opts.width), 10))
	}
	if opts.height > 0 {
		args = append(args, "--height", strconv.FormatUint(uint64(opts.height), 10))
	}
	if opts.okLabel != nil {
		args = append(args, "--ok-label", *opts.okLabel)
	}
	if opts.cancelLabel != nil {
		args = append(args, "--cancel-label", *opts.cancelLabel)
	}
	if opts.extraButton != nil {
		args = append(args, "--extra-button", *opts.extraButton)
	}
	if opts.noWrap {
		args = append(args, "--no-wrap")
	}
	if opts.ellipsize {
		args = append(args, "--ellipsize")
	}
	if opts.defaultCancel {
		args = append(args, "--default-cancel")
	}
	switch opts.icon {
	case NoIcon:
		args = append(args, "--icon-name=")
	case ErrorIcon:
		args = append(args, "--window-icon=error", "--icon-name=dialog-error")
	case WarningIcon:
		args = append(args, "--window-icon=warning", "--icon-name=dialog-warning")
	case InfoIcon:
		args = append(args, "--window-icon=info", "--icon-name=dialog-information")
	case QuestionIcon:
		args = append(args, "--window-icon=question", "--icon-name=dialog-question")
	}

	out, err := zenutil.Run(opts.ctx, args)
	if len(out) > 0 && opts.extraButton != nil &&
		string(out[:len(out)-1]) == *opts.extraButton {
		return false, ErrExtraButton
	}
	if err, ok := err.(*exec.ExitError); ok && err.ExitCode() == 1 {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err
}
