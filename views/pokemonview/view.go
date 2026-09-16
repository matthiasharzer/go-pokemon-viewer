package pokemonview

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/matthiasharzer/go-pokemon-viewer/domain/pokemon"
	"github.com/matthiasharzer/go-pokemon-viewer/logging"
	"github.com/matthiasharzer/go-pokemon-viewer/view"
	"github.com/matthiasharzer/go-pokemon-viewer/view/inmemory"
)

//go:embed pokedex.json
var pokedexRawData []byte

func toDomainTranslation(translation PokedexTranslation) pokemon.Translation {
	return pokemon.Translation{
		English: translation.English,
		German:  translation.German,
	}
}

func toDomainPokemon(pokedexPokemon PokedexPokemon) pokemon.Pokemon {
	pkm := pokemon.Pokemon{
		ID:         pokedexPokemon.ID,
		DexNr:      pokedexPokemon.DexNr,
		Generation: pokedexPokemon.Generation,
		Names:      toDomainTranslation(pokedexPokemon.Names),
		Stats: pokemon.Stats{
			Attack:  pokedexPokemon.Stats.Attack,
			Defense: pokedexPokemon.Stats.Defense,
			Stamina: pokedexPokemon.Stats.Stamina,
		},
		PrimaryType: pokemon.Type{
			Type:  pokedexPokemon.PrimaryType.Type,
			Names: toDomainTranslation(pokedexPokemon.PrimaryType.Names),
		},
		Assets: pokemon.Assets{
			Image:      pokedexPokemon.Assets.Image,
			ShinyImage: pokedexPokemon.Assets.ShinyImage,
		},
	}
	if pokedexPokemon.SecondaryType != nil {
		pkm.SecondaryType = &pokemon.Type{
			Type:  pokedexPokemon.SecondaryType.Type,
			Names: toDomainTranslation(pokedexPokemon.SecondaryType.Names),
		}
	}

	return pkm
}

func New() (view.ReadOnlyView[pokemon.Pokemon], error) {
	pokemonView := inmemory.NewView[pokemon.Pokemon]()

	var pokedex []PokedexPokemon
	err := json.Unmarshal(pokedexRawData, &pokedex)
	if err != nil {
		return nil, err
	}

	var allPokemon []pokemon.Pokemon
	for _, pkm := range pokedex {
		allPokemon = append(allPokemon, toDomainPokemon(pkm))
	}
	logging.Info(fmt.Sprintf("Loaded %d Pokémon from pokedex.json", len(allPokemon)))
	err = pokemonView.Insert(allPokemon...)
	if err != nil {
		return nil, err
	}

	return pokemonView, nil
}
