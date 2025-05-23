# krekon api

An api to get recon data from mongodb.


## CHECK LIST

+ [ ] API
    + [X] GET input entries, regex search 
    + [X] GET hosts entries, regex search 
    + [X] GET dns entries, regex search 
    + [X] GET tls entries, regex search 
    + [X] GET http entries, regex search 
    + [X] GET single http entry, based on host
    + [X] GET single dns entry, based on hostname
    + [X] GET single entry, based on hostname
    + [X] GET single tls entry, based on hostname
    + [X] GET Single host entry, based on hostname
    + [X] GET Detailed info of a host
    + [X] POST input entries 
    + [X] POST dns entries 
    + [X] POST tls entries 
    + [X] POST http entries
    + [ ] Implement pagination for get requests - 25,50,100
    + [ ] Handle errors - return json string with error instead of a text document
    + [ ] Add redis caching - set cache logic in krekon and only get cache if there is any in the API
    + [ ] Create unit and integration tests

