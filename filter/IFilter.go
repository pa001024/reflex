package filter

import ()

type IFilter interface {
	Process(src string) (dst string)
}
