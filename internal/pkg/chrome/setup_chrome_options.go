package chrome

import (
	"alice088/sparser/internal/pkg/samokat"
	"github.com/chromedp/chromedp"
	"os"
)

func SetupChromedpOptions() []chromedp.ExecAllocatorOption {
	return append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(os.Getenv("EXEC_CHROME_PATH")),
		chromedp.UserAgent(samokat.USER_AGENT),
		chromedp.Flag("headless", true),
		chromedp.DisableGPU,
		chromedp.NoSandbox,
	)
}
