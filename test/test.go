package test

/*
#include "stdio.h"
void hello(){printf("hello from test");}
*/
import "C"

func Hello() {
	C.hello()
}
