from pathlib import Path
import pygame

def ensure_game_assets_exists() -> None:

    """Ensures that game assets are downloaded since we use non-redistributable game assets"""

    assets_dir = Path("assets")
    assets = [
        "SpaceShooterAssetPack_BackGrounds.png",
        "SpaceShooterAssetPack_Characters.png",
        "SpaceShooterAssetPack_Ships.png",
        "SpaceShooterAssetPack_Projectiles.png",
        "SpaceShooterAssetPack_Miscellaneous.png",
        "SpaceShooterAssetPack_IU.png",
    ]


    for asset in assets:
        if not assets_dir.joinpath(asset).exists():
            print(f"ERROR: {asset} is not in the assets/ directory. Follow the instructions on the `README.md`  to resolve this.")
            exit(1)

def main() -> None:
    ensure_game_assets_exists()

    pygame.init()
    screen = pygame.display.set_mode((1280, 720))
    clock = pygame.time.Clock()
    running = True

    while running:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False

        screen.fill("purple")

        pygame.display.flip()

        clock.tick(60)  # limits FPS to 60

    pygame.quit()

