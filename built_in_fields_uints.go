package croconf //nolint:dupl

import "strconv"

func uintValHelper(sources []UintValueBinder, bitSize int, saveToDest func(uint64)) func(sourceNum int) Binding {
	return func(sourceNum int) Binding {
		var val uint64
		binding := sources[sourceNum].BindUintValueTo(&val)

		return wrapBinding(binding, func() error {
			if err := binding.Apply(); err != nil {
				return err
			}
			if err := checkUintBitsize(val, bitSize); err != nil {
				return err
			}
			saveToDest(val)
			return nil
		})
	}
}

func NewUintField(dest *uint, sources ...UintValueBinder) Field {
	return newField(dest, len(sources), uintValHelper(sources, strconv.IntSize, func(val uint64) {
		*dest = uint(val) // this is safe, uintValHelper checks val against bitSize
	}))
}

func NewUint8Field(dest *uint8, sources ...UintValueBinder) Field {
	return newField(dest, len(sources), uintValHelper(sources, 8, func(val uint64) {
		*dest = uint8(val) // this is safe, uintValHelper checks val against bitSize
	}))
}

func NewUint16Field(dest *uint16, sources ...UintValueBinder) Field {
	return newField(dest, len(sources), uintValHelper(sources, 16, func(val uint64) {
		*dest = uint16(val) // this is safe, uintValHelper checks val against bitSize
	}))
}

func NewUint32Field(dest *uint32, sources ...UintValueBinder) Field {
	return newField(dest, len(sources), uintValHelper(sources, 32, func(val uint64) {
		*dest = uint32(val) // this is safe, uintValHelper checks val against bitSize
	}))
}

func NewUint64Field(dest *uint64, sources ...UintValueBinder) Field {
	return newField(dest, len(sources), func(sourceNum int) Binding {
		return sources[sourceNum].BindUintValueTo(dest)
	})
}

func uintSliceHandler(newTypedSlice func(int) (add func(uint64) error, save func())) arrayHandler {
	return func(arrLength int, getElement func(int) LazySingleValueBinder) error {
		add, save := newTypedSlice(arrLength)
		for i := 0; i < arrLength; i++ {
			var val uint64
			elBinding := getElement(i).BindUintValueTo(&val)
			if err := elBinding.Apply(); err != nil {
				return err
			}
			if err := add(val); err != nil {
				return err
			}
		}
		save()
		return nil
	}
}

func NewUintSliceField(dest *[]uint, sources ...ArrayValueBinder) Field {
	return newArrayField(dest, sources, uintSliceHandler(func(arrLength int) (func(uint64) error, func()) {
		newArr := make([]uint, 0, arrLength)
		add := func(val uint64) error {
			if err := checkUintBitsize(val, strconv.IntSize); err != nil {
				return err
			}
			newArr = append(newArr, uint(val)) // this is safe
			return nil
		}
		save := func() { *dest = newArr }
		return add, save
	}))
}

func NewUint8SliceField(dest *[]uint8, sources ...ArrayValueBinder) Field {
	return newArrayField(dest, sources, uintSliceHandler(func(arrLength int) (func(uint64) error, func()) {
		newArr := make([]uint8, 0, arrLength)
		add := func(val uint64) error {
			if err := checkUintBitsize(val, 8); err != nil {
				return err
			}
			newArr = append(newArr, uint8(val)) // this is safe
			return nil
		}
		save := func() { *dest = newArr }
		return add, save
	}))
}

func NewUint16SliceField(dest *[]uint16, sources ...ArrayValueBinder) Field {
	return newArrayField(dest, sources, uintSliceHandler(func(arrLength int) (func(uint64) error, func()) {
		newArr := make([]uint16, 0, arrLength)
		add := func(val uint64) error {
			if err := checkUintBitsize(val, 16); err != nil {
				return err
			}
			newArr = append(newArr, uint16(val)) // this is safe
			return nil
		}
		save := func() { *dest = newArr }
		return add, save
	}))
}

func NewUint32SliceField(dest *[]uint32, sources ...ArrayValueBinder) Field {
	return newArrayField(dest, sources, uintSliceHandler(func(arrLength int) (func(uint64) error, func()) {
		newArr := make([]uint32, 0, arrLength)
		add := func(val uint64) error {
			if err := checkUintBitsize(val, 32); err != nil {
				return err
			}
			newArr = append(newArr, uint32(val)) // this is safe
			return nil
		}
		save := func() { *dest = newArr }
		return add, save
	}))
}

func NewUint64SliceField(dest *[]uint64, sources ...ArrayValueBinder) Field {
	return newArrayField(dest, sources, uintSliceHandler(func(arrLength int) (func(uint64) error, func()) {
		newArr := make([]uint64, 0, arrLength)
		add := func(val uint64) error {
			newArr = append(newArr, val)
			return nil
		}
		save := func() { *dest = newArr }
		return add, save
	}))
}
