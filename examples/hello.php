<?php

// test run function
foo_hello();
foo_hello3(12, "abc", 3.5);

// test run class method
$d = new PGDemo();
$r = $d->Hello1();
var_dump($r);
$r2 = $d->Hello3();
print_r($r2);
$r3 = $d->Hello4();
print_r($r3);
