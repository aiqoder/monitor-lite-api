import { ref } from "vue";

export function usePagination(_size = 10) {
    const current = ref(1)
    const size = ref(_size)
    const total = ref(0)

    function changeSize(_size: number) {
        size.value = _size
    }

    function changeCurrent(_current: number) {
        if (_current <= 1) {
            current.value = 1
            return
        }
        current.value = _current
    }

    return {
        current,
        size,
        total,
        changeSize,
        changeCurrent,
    }
}