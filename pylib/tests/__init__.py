import unittest

import os
import sys

# this makes the src directory available to the tests
sys.path.insert(0, os.path.abspath( os.path.join(os.path.dirname(__file__), '../src/') ))

if __name__ == '__main__':
    unittest.main()