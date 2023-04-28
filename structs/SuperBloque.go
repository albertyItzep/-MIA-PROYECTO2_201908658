package structs

type SuperBlock struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_blocks_count int32
	S_free_inodes_count int32
	S_mtime             [8]byte
	S_mnt_count         int32
	S_magic             int32
	S_inode_size        int32
	S_block_size        int32
	S_firts_ino         int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

type InodeTable struct {
	I_uid   int32
	I_gid   int32
	I_size  int32
	I_atime [8]byte
	I_ctime [8]byte
	I_mtime [8]byte
	I_block [16]int32
	I_type  byte
	I_perm  int32
}
type DirBlock struct {
	B_Content [4]B_content
}

type B_content struct {
	B_name  [12]byte
	B_inodp int32
}

type FileBlock struct {
	B_content [64]byte
}
