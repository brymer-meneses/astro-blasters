from space_shooter.assets import AssetLoader

import pygame

SCREEN_WIDTH = 1280
SCREEN_HEIGHT = 720

def render_background(screen: pygame.SurfaceType, asset_loader: AssetLoader) -> None:

    for x in range(0, SCREEN_WIDTH, 128):
        for y in range(0, SCREEN_HEIGHT, 256):
            index = (x + y) % 3
            screen.blit(asset_loader.backgrounds[index], (x, y))
        
    return

def main() -> None:
    asset_loader = AssetLoader()

    pygame.init()
    screen = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))
    clock = pygame.time.Clock()

    screen.fill("black")
    render_background(screen ,asset_loader)

    running = True
    while running:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False

        pygame.display.flip()

        clock.tick(60)  # limits FPS to 60

    pygame.quit()

    return

