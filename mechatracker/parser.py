import logging

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
            logging.error(f"Error while loading file at {path}: {e}")
            raise ValueError(f"File at {path} is not valid .grbr file.")
        except FileNotFoundError:
            logging.error(f"File at {path} not found.")
            raise FileNotFoundError(f"File at {path} not found.")
