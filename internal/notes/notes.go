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
	GetNoteByID(int) (*Note, error)
	GetNoteByTitle(string) (*Note, error)
}

type NoteWriter interface {
	InsertNote(Note) (int, error)
	UpdateNote(int, Note) (int64, error)
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

func (c *Client) GetByID(id int) (*Note, error) {
	return c.nr.GetNoteByID(id)
}

func (c *Client) GetByTitle(title string) (*Note, error) {
	return c.nr.GetNoteByTitle(title)
}

func (c *Client) List() ([]Note, error) {
	return c.nr.ListNotes()
}

func (c *Client) Create(n Note) (int, error) {
	return c.nw.InsertNote(n)
}

func (c *Client) Update(id int, n Note) (int64, error) {
	return c.nw.UpdateNote(id, n)
}

func (c *Client) Delete(id int) error {
	return c.nw.DeleteNote(id)
}
