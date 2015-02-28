codepen
===============


Simple API wrapper for codepen. This is useful for collection view, heart, and comment stats for your pens and collections. [http://codepen-api.appspot.com/](http://codepen-api.appspot.com/).


### Pen

specify the user and codepen ID. `http://codepen-api.appspot.com/{user}/details/{id}`

    http://codepen-api.appspot.com/drewwells/details/wBdBOB

Response

    {"comments":0,"hearts":0,"referrer":"http://codepen.io/drewwells/details/wBdBOB","views":35}

### Collection

specify the user and collection ID. `http://codepen-api.appspot.com/collection/{id}`

    http://codepen-api.appspot.com/collection/DbNZQJ/

Response

    [
      {
        "pen":"http://codepen.io/drewwells/details/qEoReo",
        "url":"qEoReo","comments":0,"views":2,"loves":0
      },
      {
        "pen":"http://codepen.io/drewwells/details/LEbPyN",
        "url":"LEbPyN","comments":0,"views":23,"loves":0},
      {
        "pen":"http://codepen.io/drewwells/details/YPpXoE",
        "url":"YPpXoE","comments":0,"views":26,"loves":0},
      {
        "pen":"http://codepen.io/drewwells/details/gbLJXR",
        "url":"gbLJXR","comments":0,"views":20,"loves":0},
      {
        "pen":"http://codepen.io/drewwells/details/KwaWva",
        "url":"KwaWva","comments":0,"views":22,"loves":0},
      {
        "pen":"http://codepen.io/drewwells/details/wBdBOB",
        "url":"wBdBOB","comments":0,"views":35,"loves":0
      }
    ]
