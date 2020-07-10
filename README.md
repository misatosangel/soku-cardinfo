# soku-cardinfo
Simple repository for access to Hisoutensoku card information

This code is really just a library around readers for Hisoutensoku's card data. For convenience, it also comes with a convenient CSV file for giving the names and costs and codes of all the cards in the game.

There is a simple reader tool for testing this and printing decks.

Example usage

`go build ./cmd/reader`

then

`./reader -f data/all_cards.csv --char Suwako -H CA 66 C9 6E`

Might give output
```
Deck for Suwako:
1 x [66] 236 Skill Card: 'Lake of Great Earth'
1 x [6e] 236 Skill Card: 'Stone Frog God'
1 x [c9] 2-Cost Spell Card: 'Spring Sign "Moriya Clear Water"'
1 x [ca] 3-Cost Spell Card: 'Party Start "2 bows, 2 Claps, 1 Bow"'
```
