import { useEffect, useRef } from 'react'

const Checkbox = ({
  indeterminate = false,
  onChange,
  checked,
}: {
  indeterminate: boolean
  onChange: any
  checked: any
}) => {
  const cRef: any = useRef()

  useEffect(() => {
    cRef.current.indeterminate = indeterminate
  }, [cRef, indeterminate])

  return (
    <input className="checkbox" type="checkbox" checked={checked} onChange={onChange} ref={cRef} />
  )
}

export default Checkbox
