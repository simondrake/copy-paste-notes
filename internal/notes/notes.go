package notes

type Client struct {
	nr NoteReader
	nw NoteWriter
}

type Note struct {
	ID              int
	Title           string
	Description     string
	CreateTimestamp string
}

type NoteReader interface {
	ListNotes() ([]Note, error)
	GetNote(int) (*Note, error)
}

type NoteWriter interface {
	InsertNote(Note) (int, error)
	UpdateNote(string, Note) (*Note, error)
	DeleteNote(int) error
}

type NoteReaderWriter interface {
	NoteReader
	NoteWriter
}

func New(nrw NoteReaderWriter) *Client {
	return &Client{
		nr: nrw,
		nw: nrw,
	}
}

func (c *Client) Get(id int) (*Note, error) {
	return c.nr.GetNote(id)
}

func (c *Client) List() ([]Note, error) {
	return c.nr.ListNotes()
}

func (c *Client) Create(n Note) (int, error) {
	return c.nw.InsertNote(n)
}

func (c *Client) Delete(id int) error {
	return c.nw.DeleteNote(id)
}
