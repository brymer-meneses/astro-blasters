from pathlib import Path
from typing import List, Tuple

import pygame
import random


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

        self.background_image = pygame.image.load(assets_dir.joinpath("SpaceShooterAssetPack_BackGrounds.png"))
        return


class Background:

    def __init__(self, 
                 screen_width: int,
                 screen_height: int,
                 asset_loader: AssetLoader) -> None:

        self.width = screen_width
        self.height = screen_height

        self.asset_loader = asset_loader
        self.positions = []
        self.surface = pygame.surface.Surface((screen_width, screen_height))

        backgrounds = []
        for (x, y) in [(0, 0), (128, 256), (256, 256)]:
            backgrounds.append(asset_loader.background_image.subsurface((x, y, 128, 256)))

        for x in range(0, screen_width, 128):
            for y in range(0, screen_height, 256):
                self.surface.blit(random.choice(backgrounds), (x, y))

        return


