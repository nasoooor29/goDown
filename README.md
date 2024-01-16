TODO: 
*   Add rolling for the log file eg. https://stackoverflow.com/questions/45440491/how-to-configure-uber-go-zap-logger-for-rolling-filesystem-log
*   Add Another Sites for manga and novels
*   Create API endpoints
*   Do the DB


CRUD mapping:
    Post    -> Create
    Get     -> Read
    Put     -> Update
    Delete  -> Delete


Api Endpoints:
    Read (for every one who have read)
        [Endpoint]/SubEndpoint        return all the media IDs as pages according to the Pagesize param 
        [Endpoint]/SubEndpoint/ID     return the requested thing
    
    MediaCollection
        Create
            Payload {
                Name            string
                createDate      int
                readingStatus   iota
                type            iota
                content         []IDs
            } -> ID
        Read
        Delete
            Goroutine running 24/7 to find the finished and remove the time based ones
        Update
            Payload {
                name            string
                readingStatus   iota
                type            iota
                content         []IDs
            } -> ID
    Media
        Create
            Payload [create by link] {
                Link                string
                LoadAllEpisodes     bool
            } -> ID
        Read
        Delete
    Episode
        Read
        Delete

    Download
        Post     Start the download 
            Payload {
                Path string
            }
        Delete   Stop the download
    status
        Read 
            get the media reading and download status
        Update {
            
        }