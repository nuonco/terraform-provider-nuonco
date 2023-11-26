package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/nuonco/nuon-go"
)

func writeDiagnosticsErr(_ context.Context, diagnostics *diag.Diagnostics, err error, op string) {
	var msg = "Please try this operation again, and if it persists, please contact us."
	userErr, ok := nuon.ToUserError(err)
	if ok {
		msg = fmt.Sprintf("Error: %s\n\nError Details: %s", userErr.Error, userErr.Description)
	}

	stateErr := &retry.UnexpectedStateError{}
	if errors.As(err, &stateErr) {
		msg = fmt.Sprintf("Error polling state change for resource. This may require manual intervention. %s\nError: %s", msg, err)
	}

	diagnostics.AddError(
		fmt.Sprintf("Unable to %s", op),
		msg,
	)
}

func logErr(ctx context.Context, err error, op string) {
	tflog.Trace(ctx, fmt.Sprintf("unable to %s: %s", op, err))
}
