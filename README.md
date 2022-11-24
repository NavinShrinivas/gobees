# gobees
Yet another Map Reduce framework, written and powered by GO!

> Output from mapper file HASSSSS to be <key,value> => comma is a must!

> NOTE : Giving custom partition function will RADICALLY slow down map reduce, this is due to the limitations of golang not havin generic implementations at runtime like Rust or Java :,(

> Note : if the  custom partition function is not following the fixed template, it may lead to infifite job.
