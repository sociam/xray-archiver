"""
Just comms step before gplaycli 
"""
import sys


def main():
    try:
        import gplaycli 
    except ImportError:
        print ("Please make sure gplaycli is installed")
    gplaycli.main()

if __name__ == '__main__':
    main()
    