<?php

// test run function
foo_hello();

// test run class method
$d = new PGDemo();
$r = $d->Hello1();
var_dump($r);
$r2 = $d->Hello3();
print_r($r2);
