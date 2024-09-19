from space_shooter.assets import AssetLoader

import pygame

def main() -> None:
    asset_loader = AssetLoader()
    screen_width = 1280
    screen_height = 720

    pygame.init()
    screen = pygame.display.set_mode((screen_width, screen_height))
    clock = pygame.time.Clock()

    running = True
    while running:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False

        screen.fill("black")
        screen.blit(asset_loader.backgrounds[0], (0, 0))

        pygame.display.flip()

        clock.tick(60)  # limits FPS to 60

    pygame.quit()

