package nft

// Metadata entity describes metadata of nft token format erc-721.
type Metadata struct {
	Description string      `json:"description"`
	ExternalURL string      `json:"external_url"`
	Image       string      `json:"image"`
	Name        string      `json:"name"`
	Attributes  []Attribute `json:"attributes"`
	// DNA is the number of the layer combination of the image, as in HashLips Art Engine metadata.
	DNA string `json:"dna,omitempty"`
}

// Attribute entity describes attributes for the item, which will show up on the OpenSea page for the item.
type Attribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}
