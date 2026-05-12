from mechatracker.parser import *
from glob import glob

import xml.etree.ElementTree as ET

def main():
    logging.basicConfig(level=logging.INFO)
    logging.info("Starting mechatracker...")
    parse_replay("ReplayFiles/2013_20260415--201415094_[Charlie]VS[BobVance].grbr")

    # for path in glob("ReplayFiles/*.grbr"):
    #     root = ET.fromstring(load_replay_bytes(path))
    #     # for child in root:
    #     #     print(child.tag, child.attrib)
    #     #     for nested_child in child:
    #     #         print("    ", nested_child.tag, nested_child.attrib)
    #     snapData = root.findall('matchDatas/MatchSnapshotData')
    #     for snap in snapData:
    #         reinforceItems = (snap.find('reinforceItems'))
    #         if reinforceItems is None:
    #             continue
    #         if len(reinforceItems) > 1:
    #             arrayOfInts = reinforceItems[1]
    #             for id in arrayOfInts:
    #                 id = id.text
    #                 if id[0] == '3':
    #                     continue
    #                 else:
    #                     print(path)
    #                     print(id)

if __name__ == "__main__":
    main()
