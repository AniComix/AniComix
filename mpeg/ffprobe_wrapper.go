package mpeg

func ffprobe() *FFCommand {
	return &FFCommand{args: []string{"ffprobe"}}
}

const (
	optShowEntries = "-show_entries"
	optOf          = "-of"
)

func (c *FFCommand) setShowEntry(entry string) *FFCommand {
	c.args = append(c.args, optShowEntries, entry)
	return c
}

func (c *FFCommand) setOf(format string) *FFCommand {
	c.args = append(c.args, optOf, format)
	return c
}
