#!/usr/bin/env node

const { Controller } = require('../dist/flip/controller');
const { Game } = require('../dist/flip/game')
const { readFileSync, writeFileSync } = require('fs');
const { join } = require('path');
const readlineSync = require('readline-sync');
const readline = require('readline');

const configPath = join(process.cwd(), 'config.json');

function waitForKeypress() {
    return new Promise((resolve) => {
        process.stdin.setRawMode(true);
        process.stdin.resume();
        process.stdin.once('data', (key) => {
            process.stdin.setRawMode(false);
            process.stdin.pause();
            resolve(key.toString());
        });
    });
}

async function main() {
    const args = process.argv.slice(2);
    const command = args[0];

    if (!command) {
        console.log('Available commands:')
        console.log('flip7 newGame: start a new game')
        console.log('flip7 join: join a game')
        return
    }

    if (command === 'newGame') {
        let num = readlineSync.question('Enter total players: ');
        const game = new Game()
        await game.init(num)
    } else if (command === 'join') {
        let name
        let secret
        let hostIp
        const controller = new Controller();
        while (true) {
            try {
                name = readlineSync.question('Enter your name: ');
                secret = readlineSync.question('Enter a password: ');
                hostIp = readlineSync.question('Enter host IP: ');

                writeFileSync(
                    configPath,
                    JSON.stringify(
                        { name, secret, hostIp, counter: -1, time: new Date().toString() },
                        null,
                        2
                    )
                );
                const ipPath = join(process.cwd(), 'ipAddress.txt');
                writeFileSync(ipPath, hostIp);
                await controller.loadMetadata()
                await controller.register(name, secret, hostIp)
                break
            } catch (error) {
                console.log(error.message)
                if (error.message == 'Room full') {
                    return
                }
            }
        }

        console.log('✔ Registered successfully!');
        const data = JSON.parse(readFileSync(configPath, 'utf8'));
        data.counter = -1;
        data.time = new Date().toString();
        writeFileSync(configPath, JSON.stringify(data, null, 2));

        await controller.init(data.name, data.secret, data.hostIp);
        await controller.refreshMetadata();
        let isRunning = false;

        setInterval(() => {
            if (!isRunning) {
                isRunning = true;
                loopThis().finally(() => {
                    isRunning = false;
                });
            }
        }, 1000);

        async function loopThis() {
            try {
                await controller.refreshMetadata();
                console.clear();
                await controller.getRoundInfo();

                const turnName = await controller.getTurn();

                let localData = await controller.getMetadata()
                if (localData.totalPlayers > localData.players.length) {
                    console.log("Waiting for players to join...")
                    return
                }

                if (turnName !== data.name) {
                    await controller.getRoundInfo(true);
                    return;
                }



                console.log(`\n🎮 Welcome ${data.name}! It's your turn.`);

                const metadata = await controller.getMetadata();
                const turnType = metadata.turnCounter.turn;

                if (turnType === 'Hit') {
                    while (true) {
                        console.log('\nPress [Space] to HIT or [f] to FREEZE');
                        const key = await waitForKeypress();

                        if (key === ' ') {
                            try {
                                await controller.refreshMetadata();
                                await controller.hit();
                                await controller.uploadMetadata();
                                break;
                            } catch (err) {
                                console.log(`❌ Error (hit): ${err.message || err}`);
                            }
                        } else if (key.toLowerCase() === 'f') {
                            try {
                                await controller.refreshMetadata();
                                await controller.freeze();
                                await controller.uploadMetadata();
                                break;
                            } catch (err) {
                                console.log(`❌ Error (freeze): ${err.message || err}`);
                            }
                        } else {
                            console.log('❗ Invalid key. Use Space for hit or f for freeze.');
                        }
                    }

                } else if (turnType === 'Hit-only') {
                    while (true) {
                        console.log('\n🔒 You can only HIT. Press [Space] to HIT');
                        const key = await waitForKeypress();

                        if (key === ' ') {
                            try {
                                await controller.refreshMetadata();
                                await controller.hit();
                                await controller.uploadMetadata();
                                break;
                            } catch (err) {
                                console.log(`❌ Error (hit-only): ${err.message || err}`);
                            }
                        } else {
                            console.log('❗ Invalid key. Only HIT is allowed (Press Space).');
                        }
                    }

                } else if (turnType === 'Choose') {
                    const players = metadata.players.filter(
                        (p) => p.status === 'active'
                    );

                    if (players.length === 0) {
                        console.log('❗ No valid players to choose.');
                        return;
                    }

                    while (true) {
                        console.log('\n Choose a player by number:');
                        players.forEach((p, index) => {
                            console.log(`[${index}] ${p.name}`);
                        });

                        const indexStr = readlineSync.question('\nEnter number: ').trim();
                        const index = parseInt(indexStr);

                        if (isNaN(index) || index < 0 || index >= players.length) {
                            console.log('❗ Invalid selection.');
                            continue;
                        }

                        const chosenName = players[index].name;

                        try {
                            await controller.refreshMetadata();
                            await controller.choose(chosenName);
                            await controller.uploadMetadata();
                            break;
                        } catch (err) {
                            console.log(`❌ Error (choose): ${err.message || err}`);
                        }
                    }
                } else {
                    console.log(`❗ Unknown turn type: ${turnType}`);
                }

            } catch (err) {
                console.log(`❌ Fatal Error in loopThis: ${err.message || err}`);
            }
        }
    } else {
        console.log(`Unknown command: ${command}`);
    }
}

main();