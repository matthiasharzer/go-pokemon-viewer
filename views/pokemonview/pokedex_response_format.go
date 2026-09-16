package pokemonview

type responsePokedexTranslation struct {
	English string `json:"English"`
	German  string `json:"German"`
}

type responsePokedexType struct {
	Type  string                     `json:"type"`
	Names responsePokedexTranslation `json:"names"`
}

type responsePokedexStats struct {
	Attack  int `json:"attack"`
	Defense int `json:"defense"`
	Stamina int `json:"stamina"`
}

type responsePokedexAssets struct {
	Image      string `json:"image"`
	ShinyImage string `json:"shinyImage"`
}

type responsePokedexPokemon struct {
	ID            string                     `json:"id"`
	DexNr         int                        `json:"dexNr"`
	Generation    int                        `json:"generation"`
	Names         responsePokedexTranslation `json:"names"`
	Stats         responsePokedexStats       `json:"stats"`
	PrimaryType   responsePokedexType        `json:"primaryType"`
	SecondaryType *responsePokedexType       `json:"secondaryType,omitempty"`
	Assets        responsePokedexAssets      `json:"assets"`
}
