package structs

type Partition struct {
	Part_status byte
	Part_type   byte
	Part_fit    byte
	Part_start  uint32
	Part_size   uint32
	Part_name   [16]byte
}
type MBR struct {
	Mbr_tamano         uint32
	Mbr_fecha_creacion [8]byte
	Mbr_dsk_signature  uint32
	Dsk_fit            byte
	Mbr_partition1     Partition
	Mbr_partition2     Partition
	Mbr_partition3     Partition
	Mbr_partition4     Partition
}

type EBR struct {
	Part_status, Part_fit byte
	Part_start, Part_size int32
	Part_next             int32
	Part_name             [16]byte
}

func (tmp *MBR) InitialPartitions() {
	tmp.Mbr_partition1.Part_status = 'f'
	tmp.Mbr_partition2.Part_status = 'f'
	tmp.Mbr_partition3.Part_status = 'f'
	tmp.Mbr_partition4.Part_status = 'f'

	tmp.Mbr_partition1.Part_type = 'f'
	tmp.Mbr_partition2.Part_type = 'f'
	tmp.Mbr_partition3.Part_type = 'f'
	tmp.Mbr_partition4.Part_type = 'f'

	tmp.Mbr_partition1.Part_fit = 'v'
	tmp.Mbr_partition2.Part_fit = 'v'
	tmp.Mbr_partition3.Part_fit = 'v'
	tmp.Mbr_partition4.Part_fit = 'v'
}
