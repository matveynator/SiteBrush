package Data

type Post struct {
  Id           int64
  OwnerId      int
  EditorId     int
  RequestUri   string
  Date         int64
  Title        string
  Body         string
  Header       string
  Tags         string
  Revision     int
  Domain       string
  Status       string
  Published    bool
}


