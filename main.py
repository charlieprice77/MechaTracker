from mechatracker.parser import *

def main():
    logging.basicConfig(level=logging.INFO)
    logging.info("Starting mechatracker...")
    load_replay_bytes("ReplayFiles/2013_20260415--201415094_[Charlie]VS[BobVance].grbr")


if __name__ == "__main__":
    main()
