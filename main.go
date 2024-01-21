package main

type KeyStore[K comparable, V any] interface {
	Set(K, V)
	Get(K) V
	Delete(K)
	Count() int
	Begin()
	End()
	Rollback()
	Commit()
}
