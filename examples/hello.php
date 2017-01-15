<?php

// test run function
function test_funcs() {
    foo_hello();
    foo_hello3(12, "abc", 3.5);
}

// test run class method
function test_classes() {
    $d = new PGDemo();
    $r = $d->Hello1();
    var_dump($r);
    $r2 = $d->Hello3();
    print_r("index array result: " . var_export($r2, true));
    $r3 = $d->Hello4();
    print_r("assoc array result: " .var_export($r3, true));

    // test for executing scope destruct and execute finished scope destruct
    if (rand() % 2 == 1) {
        $d = null;
    }
}

// test global ini vars
function test_inis() {
    $inis = ['phpgo.hehe_int', 'phpgo.hehe_bool', 'phpgo.hehe_long', 'phpgo.hehe_string'];
    for ($idx = 0; $idx < count($inis); $idx ++) {
        print($inis[$idx] . ' = ' . var_export(ini_get($inis[$idx]), true) . "\n");
    }

    for ($idx = 0; $idx < count($inis); $idx ++) {
        $ov = ini_get($inis[$idx]);
        ini_set($inis[$idx], $ov . "++");
        print($inis[$idx] . ' = ' . var_export(ini_get($inis[$idx]), true) . "\n");
    }

    $inis = ['pg0.hehe', 'pg0.hehe2'];
    for ($idx = 0; $idx < count($inis); $idx ++) {
        print($inis[$idx] . ' = ' . var_export(ini_get($inis[$idx]), true) . "\n");
    }

    for ($idx = 0; $idx < count($inis); $idx ++) {
        $ov = ini_get($inis[$idx]);
        ini_set($inis[$idx], $ov . "++");
        print($inis[$idx] . ' = ' . var_export(ini_get($inis[$idx]), true) . "\n");
    }

    foreach (ini_get_all() as $k => $v) {
        if (substr($k, 0, 6) == "phpgo.") {
            print($k . ' = ' . str_replace("\n", "", var_export($v, true)) . "\n");
        } else {
            // print($k . ' = ' . var_export($v, true) . "\n");
        }
    }
}

test_funcs();
test_classes();
test_inis();

