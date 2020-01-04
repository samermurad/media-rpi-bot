package cmds

type Command interface {
	Exec(data interface{}) error
	Args() map[string]interface{}
}
