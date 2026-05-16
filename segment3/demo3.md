## Pointers and structs

You are building a simple player management system for a game. A player has a name, a health value (0–100), a score, and a level.
Write a program that models a single player and simulates a short sequence of events during a game session.

Your Player should support the following behaviours:

- Taking damage — health cannot go below 0
- Healing — health cannot go above 100
- Gaining score points
- Levelling up — increases the player's level by 1 and fully restores their health to 100
- Describing themselves — returns a formatted string showing all of their current stats

In main, create a player and run them through the following sequence of events:

- Print their starting stats
- They take 40 damage and gain 150 points
- Print their stats
- They heal 20 and gain 200 points
- Print their stats
- They take 999 damage (the game doesn't have to end here - we are just testing the Player struct)
- Print their stats
- They level up
- Print their final stats

