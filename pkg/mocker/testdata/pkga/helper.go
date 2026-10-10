package pkga

type A1 struct{}

type CaseA0 interface{ MethodA0() }

type CaseA1 interface{ MethodA1(a A1) A1 }
