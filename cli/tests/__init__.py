# Copyright (c) Microsoft Corporation.
# Licensed under the MIT license.
import logging
import os
from pathlib import Path
import unittest

class TestCaseBase(unittest.TestCase):

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.logger = logging.getLogger(self.__class__.__name__)
        
        dir_path = os.path.dirname(os.path.realpath(__file__))

        if dir_path.endswith('/tests'):
            self.data_path = Path(os.path.join(dir_path, 'testdata')).resolve()
        else:
            self.data_path = Path(os.path.join(dir_path, '../testdata')).resolve()
