import logging
from dataclasses import dataclass
from enum import Enum
import xml.etree.ElementTree as ET

class GameMode(Enum):
    """Represents the game mode of a BattleRecord."""
    FOUR_PLAYER_FFA = 'VS_4_Scuffle'
    ONE_VS_ONE = 'VS_1_1'
    TWO_VS_TWO = 'VS_2_2'

@dataclass
class UnitOffering:
    """Contains a single unit offering for a unit drop."""
    unitId: int
    unitRank: int
    numPacks: int
    cost: int | None

@dataclass
class CardSelection:
    """Containts 4 cards that are offerend to a player for a round."""
    cardIds: list[int]
    selectedCardId: int | None
    skipped: bool
    skipValue: int

@dataclass
class UnitDrop:
    """Contains the 4 offered unit drops for a round."""
    offerings: list[UnitOffering]
    selectedUnitId: int | None
    skipped: bool
    skipValue: int
    canBeSkipped: bool

@dataclass
class RoundRecord:
    """Represents a single round from a game record."""
    roundNumber: int
    unitDrop: UnitDrop
    cardSelection: CardSelection

@dataclass
class PlayerRecord:
    playerName: str
    playerTeam: int
    roundRecords: list[RoundRecord]

@dataclass
class BattleRecord:
    """Represents a BattleRecord object from a .grbr file."""
    playerRecords: list[PlayerRecord]
    gameMode: GameMode

def load_replay_bytes(path: str) -> bytes:
    """Loads a .grbr file from the given path and returns the bytes of the XML section."""
    final_xml_bytes = b"</BattleRecord>"
    with open(path, "rb") as f:
        try:
            f_bytes = f.read()
            xml_start_index = f_bytes.index(b"<?xml")
            xml_end_index = f_bytes.index(final_xml_bytes, xml_start_index) + len(final_xml_bytes)
            logging.info(f"Found XML start at {xml_start_index} and end at {xml_end_index}")
            return f_bytes[xml_start_index:xml_end_index]
        except ValueError:
            logging.error(f"Error while loading file at {path}")
            raise ValueError(f"File at {path} is not valid .grbr file.")
        except FileNotFoundError:
            logging.error(f"File at {path} not found.")
            raise FileNotFoundError(f"File at {path} not found.")

def parse_unit_offering(drop_id: int) -> UnitOffering:
    """Parses a unit offering from a given drop ID and returns a UnitOffering object."""

    return UnitOffering(unitId=0, unitRank=0, numPacks=0, cost=None)

def parse_replay(path: str) -> BattleRecord:
    """Parses a .grbr file from the given path and returns a BattleRecord object."""
    xml_bytes = load_replay_bytes(path)
    root = ET.fromstring(xml_bytes)
    # for child in root:
    #     print(child.tag, child.attrib)
    #     for nested_child in child:
    #         print("    ", nested_child.tag, nested_child.attrib)

    gameMode = GameMode(root.find("BattleInfo/MatchMode").text)
    print(gameMode)
    raise NotImplementedError
