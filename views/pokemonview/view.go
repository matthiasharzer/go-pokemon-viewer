package pokemonview

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/matthiasharzer/go-pokemon-viewer/domain/pokemon"
	"github.com/matthiasharzer/go-pokemon-viewer/logging"
	"github.com/matthiasharzer/go-pokemon-viewer/view"
	"github.com/matthiasharzer/go-pokemon-viewer/view/inmemory"
)

const pokedexURL = "https://pokemon-go-api.github.io/pokemon-go-api/api/pokedex.json"

func toDomainTranslation(translation responsePokedexTranslation) pokemon.Translation {
	return pokemon.Translation{
		English: translation.English,
		German:  translation.German,
	}
}

func toDomainPokemon(pokedexPokemon responsePokedexPokemon) pokemon.Pokemon {
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

func fetchPokedex() ([]pokemon.Pokemon, error) {
	response, err := http.Get(pokedexURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	var pokedex []responsePokedexPokemon
	err = json.NewDecoder(response.Body).Decode(&pokedex)
	if err != nil {
		return nil, err
	}

	var allPokemon []pokemon.Pokemon
	for _, pkm := range pokedex {
		allPokemon = append(allPokemon, toDomainPokemon(pkm))
	}
	return allPokemon, nil
}

func refetchPokedex(view view.View[pokemon.Pokemon]) error {
	pokedex, err := fetchPokedex()
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("fetched %d Pokémon from pokedex URL", len(pokedex)))
	err = view.ReplaceAll(pokedex...)
	if err != nil {
		return err
	}
	return nil
}

func fetchRoutine(ctx context.Context, interval time.Duration, view view.View[pokemon.Pokemon]) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logging.Info("updating pokedex")
			err := refetchPokedex(view)
			if err != nil {
				logging.Warn("failed to fetch pokedex", "error", err)
			}
		}
	}
}

func New(ctx context.Context, fetchInterval time.Duration) (view.ReadOnlyView[pokemon.Pokemon], error) {
	pokemonView := inmemory.NewView[pokemon.Pokemon]()

	err := refetchPokedex(pokemonView)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch initial pokedex: %w", err)
	}

	if fetchInterval > 0 {
		go fetchRoutine(ctx, fetchInterval, pokemonView)
	}

	return pokemonView, nil
}
