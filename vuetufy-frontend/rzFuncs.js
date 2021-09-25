    function myFetch() {
        console.log("doing fetch")
        fetch('http://localhost:8080/albums')
        .then(response => response.text())
        .then(data =>  myv.output = JSON.parse(data))
    }


    function myDeleteFromSQLDatabase(theItem) {
        console.log("deleting " + theItem)
        fetch('http://localhost:8080/deletealbum/' + theItem, {
            method: 'DELETE',
            headers: {
            'Content-type': 'application/json'
            }})
        .then(response => response.text())
        // .then(data =>  JSON.parse(data))
    }

    function myEditAlbumInSQLDatabase(theAlbumID, theTitle, theArtist, thePrice) {
        console.log("editing " + theAlbumID)
        fetch('http://localhost:8080/editalbum?id=' + theAlbumID + "&title=" + theTitle + "&artist=" + theArtist + "&price=" + thePrice)
        .then(response => response.text())
        // .then(data =>  JSON.parse(data))
    }


    async function myAdd(data) {
        await fetch('http://localhost:8080/addalbum', {
            method: 'POST',
            headers: {
            'Content-type': 'application/json'
            },
            body: data})
        .then(response => response.text())
        .then(data =>  JSON.parse(data))
    }
