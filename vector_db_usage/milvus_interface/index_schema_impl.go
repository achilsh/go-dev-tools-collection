package mivus_interface

const (
	IndexTypeAutoIndex          = 0  // autoIndex
	IndexTypediskANNIndex       = 1  // diskANNIndex
	IndexTypeflatIndex          = 2  // flatIndex
	IndexTypebinFlatIndex       = 3  //binFlatIndex
	IndexTypegpuBruteForceIndex = 4  //gpuBruteForceIndex
	IndexTypegpuIVFFlatIndex    = 5  //gpuIVFFlatIndex
	IndexTypegpuIVFPQIndex      = 6  //gpuIVFPQIndex
	IndexTypegpuCagra           = 7  //gpuCagra
	IndexTypehnswIndex          = 8  //hnswIndex
	IndexTypeGenericIndex       = 9  //GenericIndex
	IndexTypeivfFlatIndex       = 10 // ivfFlatIndex
	IndexTypeivfPQIndex         = 11 // ivfPQIndex

	IndexTypeivfSQ8Index = 12 // ivfSQ8Index
	IndexTypebinIvfFlat  = 13 // binIvfFlat

	IndexTypeJSONPathIndex = 14 // JSONPathIndex

	// 标量说索引类型
	IndexTypeTrieIndex     = 15 // TrieIndex
	IndexTypeInvertedIndex = 16 // InvertedIndex
	IndexTypeSortedIndex   = 17 //SortedIndex
	IndexTypeBitmapIndex   = 18 // BitmapIndex

	//
	IndexTypeSCANNIndex          = 30 // SCANNIndex
	IndexTypesparseInvertedIndex = 31 // sparseInvertedIndex
	IndexTypesparseWANDIndex     = 32 // sparseWANDIndex

	IndexTypesparseAnnParam = 33 // sparseAnnParam

)
