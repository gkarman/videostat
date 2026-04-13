package command

type UI struct{}

func NewUI() *UI {
	return &UI{}
}

func (u *UI) CommandsText() string {
	return `<b>Доступные команды</b>

🚀 <code>/start</code>
📖 <code>/help</code>
➕ <code>/create_blogger</code>
📋 <code>/list_bloggers</code>`
}