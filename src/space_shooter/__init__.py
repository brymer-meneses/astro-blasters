from space_shooter.assets import AssetLoader, Background
import pygame

SCREEN_WIDTH = 1280
SCREEN_HEIGHT = 720
SCROLL_SPEED = 1

def main() -> None:
    asset_loader = AssetLoader()

    pygame.init()
    screen = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))
    clock = pygame.time.Clock()

    bg = Background(SCREEN_WIDTH, SCREEN_HEIGHT, asset_loader)
    scroll = 0

    running = True
    while running:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False

        scroll %= bg.surface.get_height()

        screen.blit(bg.surface, (0, scroll))
        screen.blit(bg.surface, (0, scroll - bg.height))

        scroll += SCROLL_SPEED

        pygame.display.flip()
        clock.tick(60)

    pygame.quit()
