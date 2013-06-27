pagerank
========

[![endorse](http://api.coderwall.com/dcadenas/endorsecount.png)](http://coderwall.com/dcadenas)

A Go language [PageRank](http://en.wikipedia.org/wiki/PageRank) implementation.

[![](http://upload.wikimedia.org/wikipedia/commons/thumb/f/fb/PageRanks-Example.svg/596px-PageRanks-Example.svg.png)](http://en.wikipedia.org/wiki/PageRank)

Description
-----------

A library to calculate the [PageRank](http://en.wikipedia.org/wiki/PageRank) of a big directed graph. It's intended to be used for big but not huge graphs, as 
those are better processed with a map-reduce distributed solution.
This is a port from ruby's [rankable_graph](http://github.com/dcadenas/rankable_graph) gem.

Usage
-----

```go
package main

import "github.com/dcadenas/pagerank"
import "fmt"

func main(){
  graph := pagerank.New()

  //First we draw our directed graph using the link method which receives as parameters two identifiers.   
  //The only restriction for the identifiers is that they should be integers.
  graph.Link(1234, 4312)
  graph.Link(9876, 4312)
  graph.Link(4312, 9876)
  graph.Link(8888, 4312)

  probability_of_following_a_link := 0.85 // The bigger the number, less probability we have to teleport to some random link
  tolerance := 0.0001 // the smaller the number, the more exact the result will be but more CPU cycles will be needed

  graph.Rank(probability_of_following_a_link, tolerance, func(identifier int, rank float64) {
    fmt.Println("Node", identifier, "rank is", rank)
  })
}
```

Which outputs

    Node 1234 rank is 0.03750000000000001
    Node 4312 rank is 0.4797515116401361
    Node 9876 rank is 0.44524848835986397
    Node 8888 rank is 0.03750000000000001

This ranks represent the probabilities that a certain node will be visited.

For more examples please refer to the [tests](https://github.com/dcadenas/pagerank/blob/master/pagerank_test.go).

Note on Patches/Pull Requests
-----------------------------

* Fork the project.
* Make your feature addition or bug fix.
* Add tests for it. This is important so I don't break it in a
  future version unintentionally.
* Commit, do not mess with rakefile, version, or history.
  (if you want to have your own version, that is fine but bump version in a commit by itself I can ignore when I pull)
* Send me a pull request. Bonus points for topic branches.

Copyright
---------

Author: [Daniel Cadenas](http://danielcadenas.com)

Copyright (c) 2013. See LICENSE for details.
