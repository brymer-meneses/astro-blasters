from pathlib import Path
import pygame


class AssetLoader:

    def __init__(self):
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
                print(f"ERROR: {asset} is not in the `assets/` directory. Follow the instructions on the `README.md`  to resolve this.")
                exit(1)

        background = pygame.image.load(assets_dir.joinpath("SpaceShooterAssetPack_BackGrounds.png"))

        self.backgrounds = []
        for (x, y) in [(0, 0), (128, 256), (256, 256)]:
            bg = background.subsurface((x, y, 128, 256))
            self.backgrounds.append(bg)

