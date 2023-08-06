# HowLongToBeat

[![License](https://img.shields.io/badge/license-BSD_3--Clause-blue.svg?style=flat-square)](LICENSE.md)
![Documentation](https://img.shields.io/badge/documentation-available-yellow.svg?style=flat-square)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat-square&logo=GitHub)](https://github.com/forbiddencoding/howlongtobeat/issues)

[HowLongToBeat.com](https://howlongtobeat.com) provides information and data about video games, and how long it will
approximately take to finish them.

This library is a simple wrapper api to fetch data from HowLongToBeat. HowLongToBeat is an awesome website and a great
service, free and also living from community data.

You get the option to get raw and extensive game data or a simplified, formatted version of it.

> &#9888; **Disclaimer:** This library is not affiliated nor endorsed with HowLongToBeat.com or Ziff Davis LLC
> in any way. Please use this library responsibly and do not abuse or overload the HowLongToBeat servers.

----

## Contents

* [Quickstart](#quickstart)
* [API](#api)
    * [Search](#search)
    * [SearchSimple](#searchsimple)
    * [Detail](#detail)
    * [DetailSimple](#detailsimple)
    * [Reduce](#reduce)
* [Similar projects in different languages](#similar-projects-in-different-languages)
* [Troubleshooting](#troubleshooting)
* [Contributing](#contributing)
* [License](#license)

## Quickstart

> &#9888; Verify you are running Golang version **1.20**, by running `go version` in a terminal or console window. Older
> versions can produce errors, but newer versions are usually fine.
>
> &#9888; HowLongToBeat uses go modules, there is no dependency on `$GOPATH` variable.

#### Add **HowLongToBeat** to your project:

```bash
go get -u github.com/forbiddencoding/howlongtobeat
```

#### Building a HowLongToBeat instance:

```go
package main

import (
	"context"
	"github.com/forbiddencoding/howlongtobeat"
)

func main() {
	hltb, err := howlongtobeat.New()
	if err != nil {
		// error handling
	}
	// ...
}
```

#### Searching for a game:

```go
// ...
searchResults, err := hltb.Search(context.TODO(), "The Witcher 3: Wild Hunt", hltb.SearchModifierNone, nil)
if err != nil {
// error handling
}
// ...
```

#### Getting details for a game by its ID:

```go
// ...
gameDetails, err := hltb.Detail(context.TODO(), 40171)
if err != nil {
// error handling
}
// ...
```

## API

> A note about the returned data:
> When using the `Simple` methods, the data is already formatted and cleaned up, reducing it to the most commonly used
> data.
> Part of the formatting includes converting the time values from seconds to hours.
>
> If you find that the `Simple` methods, or the `Reduce` method respectively, trim away too much data, please open an
> issue and let me know. I am happy to add more data to the `Simple` methods, if it is required by the community.

### Search

Search allows you to search for games by their name. It returns a raw list of search results.
The result usually contains most of the data you would want or need, but might need to be formatted or cleaned up.

##### Parameters

| Name         | Type                    | Description                                                                                                                                                                                                                                                                                                     |
|--------------|-------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `query`      | `string`                | The title of the game or DLC.                                                                                                                                                                                                                                                                                   |
| `modifier`   | `SearchModifier`        | The search modifier to use. Possible values are:<br/> `SearchModifierNone` for the default behaviour returning both games as well as DLCs the default behaviour<br/>`SearchModifierOnlyDLC` to only get DLCs matching the search term or compatible with the game<br/> `SearchModiferHideDLC` to only get games |
| `pagination` | `*SearchGamePagination` | Used for custom pages sizes or pagination if too many matching result have been found by the HLTB API.<br/>The default page size is 20.                                                                                                                                                                         |

#### Usage

```go
// ...
searchResults, err := hltb.Search(context.TODO(), "The Witcher 3: Wild Hunt", hltb.SearchModifierNone, nil)
if err != nil {
// error handling
}
// ...
```

##### Returns

````json
{
  "color": "blue",
  "title": "",
  "category": "games",
  "count": 9,
  "page_current": 0,
  "page_total": 0,
  "page_size": 0,
  "data": [
    {
      "count": 9,
      "game_id": 10270,
      "game_name": "The Witcher 3: Wild Hunt",
      "game_name_date": 0,
      "game_alias": "The Witcher III: Wild Hunt",
      "game_type": "game",
      "game_image": "10270_The_Witcher_3_Wild_Hunt.jpg",
      "comp_lvl_combine": 0,
      "comp_lvl_sp": 1,
      "comp_lvl_co": 0,
      "comp_lvl_mp": 0,
      "comp_lvl_spd": 1,
      "comp_main": 184978,
      "comp_plus": 371370,
      "comp_100": 624261,
      "comp_all": 369324,
      "comp_main_count": 2316,
      "comp_plus_count": 5663,
      "comp_100_count": 1939,
      "comp_all_count": 9918,
      "invested_co": 0,
      "invested_mp": 0,
      "invested_co_count": 0,
      "invested_mp_count": 0,
      "count_comp": 18133,
      "count_speedrun": 16,
      "count_backlog": 18719,
      "count_review": 4609,
      "review_score": 94,
      "count_playing": 426,
      "count_retired": 1059,
      "profile_dev": "CD Projekt RED",
      "profile_popular": 1805,
      "profile_steam": 292030,
      "profile_platform": "Nintendo Switch, PC, PlayStation 4, PlayStation 5, Xbox One, Xbox Series X/S",
      "release_world": 2015,
      "similarity": 0.33
    }
    // ...
  ]
}
````

### SearchSimple

SearchSimple allows you to search for games by their name. It returns a list of search results with a simplified
structure, omitting data which is not usually needed.

##### Parameters

| Name       | Type             | Description                                                                                                                                                                                                                                                                               |
|------------|------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `query`    | `string`         | The title of the game or DLC.                                                                                                                                                                                                                                                             |
| `modifier` | `SearchModifier` | The search modifier to use. Possible values are:<br/> `SearchModifierNone` for the default behaviour returning both games as well as DLCs<br/>`SearchModifierOnlyDLC` to only get DLCs matching the search term or compatible with the game<br/> `SearchModiferHideDLC` to only get games |

#### Usage

```go
// ...
searchResults, err := hltb.SearchSimple(context.TODO(), "The Witcher 3", SearchModifierNone)
if err != nil {
// error handling
}
// ...
```

##### Returns

````json
[
  {
    "game_id": 10270,
    "game_name": "The Witcher 3: Wild Hunt",
    "profile_platform": "Nintendo Switch, PC, PlayStation 4, PlayStation 5, Xbox One, Xbox Series X/S",
    "game_image": "10270_The_Witcher_3_Wild_Hunt.jpg",
    "comp_main": 51,
    "comp_plus": 103,
    "comp_all": 103,
    "similarity": 0.33
  }
  // ...
]
````

### Detail

Detail allows you to get raw detailed information about a game or DLC by its ID.

##### Parameters

| Name     | Type  | Description                                                                                                                                       |
|----------|-------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| `gameID` | `int` | The ID of the game or DLC. Can be found as `GameID` by the search result or in the search bar of your browser when visiting the howlongtobeat.com |

#### Usage

```go
// ...
game, err := hltb.Detail(context.TODO(), 10270)
if err != nil {
// error handling
}
// ...
```

#### Returns

See response in example [here](examples/detail.json).

### DetailSimple

DetailSimple allows you to get detailed information about a game or DLC by its ID. It returns a simplified structure (
similar to `SearchSimple`), omitting data which is not usually needed.

#### Parameters

| Name     | Type  | Description                                                                                                                                       |
|----------|-------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| `gameID` | `int` | The ID of the game or DLC. Can be found as `GameID` by the search result or in the search bar of your browser when visiting the howlongtobeat.com |

#### Usage

```go
// ...
game, err := hltb.DetailSimple(context.TODO(), 10270)
if err != nil {
// error handling
}
// ...
```

#### Returns

````json
{
  "game_id": 10270,
  "game_name": "The Witcher 3: Wild Hunt",
  "profile_platform": "Nintendo Switch, PC, PlayStation 4, PlayStation 5, Xbox One, Xbox Series X/S",
  "game_image": "10270_The_Witcher_3_Wild_Hunt.jpg",
  "comp_main": 51,
  "comp_plus": 103,
  "comp_all": 103
}
````

### Reduce

If you used the `Search` or `Detail` function, you can call `Reduce` on the result variable to reduce the data in the
same way as `SearchSimple` and`DetailSimple` would do.

```go
// ...
game, err := hltb.Detail(context.TODO(), 10270)
if err != nil {
// error handling
}

simpleGame := game.Reduce()
// ...
```

## Similar projects in different languages

| Project                                                                                         | Language   |
|-------------------------------------------------------------------------------------------------|------------|
| [ckatzorke/howlongtobeat](https://github.com/ckatzorke/howlongtobeat)                           | TypeScript |
| [ScrappyCocco/HowLongToBeat-PythonAPI](https://github.com/ScrappyCocco/HowLongToBeat-PythonAPI) | Python     |
| [ivankayzer/howlongtobeat](https://github.com/ivankayzer/howlongtobeat)                         | PHP        |
| [jameslieu/howlongtobeat](https://github.com/jameslieu/howlongtobeat)                           | C#         |
| [saturnavt/howlongtobeat-api](https://github.com/saturnavt/howlongtobeat-api)                   | Rust       |

## Troubleshooting

For help with common problems, please see [TROUBLESHOOTING.md](TROUBLESHOOTING.md). If you encounter any other issues,
feel free to open an issue.

## Contributing

We'd love your help in making this library better. Please review our [contribution guidelines](CONTRIBUTING.md).

## License

This project is licensed under the terms of the BSD-3 License. See the [LICENSE](LICENSE.md) for more details.