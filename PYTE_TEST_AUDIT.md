# Pyte test coverage audit

Total pyte tests: 105
Mapped via comments: 90
Missing: 28

## Missing tests
- test_history.py::test_cursor_is_hidden
- test_input_output.py::test_input_output
- test_screen.py::test_clear_tabstops
- test_screen.py::test_cursor_back
- test_screen.py::test_cursor_down
- test_screen.py::test_cursor_forward
- test_screen.py::test_cursor_up
- test_screen.py::test_display_wcwidth
- test_screen.py::test_initialize_char
- test_screen.py::test_private_sgr_sequence_ignored
- test_screen.py::test_remove_non_existant_attribute
- test_screen.py::test_reset_works_between_attributes
- test_screen.py::test_restore_cursor_with_none_saved
- test_screen.py::test_save_cursor
- test_screen.py::test_tabstops
- test_stream.py::test_compatibility_api
- test_stream.py::test_control_characters
- test_stream.py::test_define_charset
- test_stream.py::test_interrupt
- test_stream.py::test_linefeed
- test_stream.py::test_missing_params
- test_stream.py::test_non_csi_sequences
- test_stream.py::test_non_utf8_shifts
- test_stream.py::test_overflow
- test_stream.py::test_reset_mode
- test_stream.py::test_set_mode
- test_stream.py::test_set_title_icon_name
- test_stream.py::test_unknown_sequences

## Mapped tests
| pyte test | go test |
| --- | --- |
| test_diff.py::test_draw_multiple_chars_wrap | TestPyteDiffDrawMultipleCharsWrap |
| test_diff.py::test_draw_wrap | TestPyteDiffDrawWrap |
| test_diff.py::test_erase_in_display | TestPyteDiffEraseInDisplay |
| test_diff.py::test_index | TestPyteDiffIndex |
| test_diff.py::test_insert_delete_lines | TestPyteDiffInsertDeleteLines |
| test_diff.py::test_mark_single_line | TestPyteDiffMarkSingleLine |
| test_diff.py::test_mark_whole_screen | TestPyteDiffMarkWholeScreen |
| test_diff.py::test_modes | TestPyteDiffModes |
| test_diff.py::test_reverse_index | TestPyteDiffReverseIndex |
| test_history.py::test_cursor_hidden | TestPyteHistoryCursorHidden |
| test_history.py::test_draw | TestPyteHistoryDraw |
| test_history.py::test_ensure_width | TestPyteHistoryEnsureWidth |
| test_history.py::test_erase_in_display | TestPyteHistoryEraseInDisplay |
| test_history.py::test_index | TestPyteHistoryIndex |
| test_history.py::test_next_page | TestPyteHistoryNextPage |
| test_history.py::test_not_enough_lines | TestPyteHistoryNotEnoughLines |
| test_history.py::test_prev_page | TestPyteHistoryPrevPageRatio |
| test_history.py::test_reverse_index | TestPyteHistoryReverseIndex |
| test_screen.py::test_alignment_display | TestPyteAlignmentDisplay |
| test_screen.py::test_attributes | TestAttributes |
| test_screen.py::test_attributes_reset | TestAttributesReset |
| test_screen.py::test_backspace | TestPyteBackspace |
| test_screen.py::test_blink | TestBlink |
| test_screen.py::test_carriage_return | TestPyteCarriageReturn |
| test_screen.py::test_clear_tab_stops | TestPyteClearTabStops |
| test_screen.py::test_colors | TestColors |
| test_screen.py::test_colors24_bit | TestPyteColors24Bit |
| test_screen.py::test_colors24bit | TestColors24Bit |
| test_screen.py::test_colors256 | TestColors256 |
| test_screen.py::test_colors256_missing_attrs | TestPyteColors256MissingAttrs |
| test_screen.py::test_colors_aixterm | TestColorsAixterm |
| test_screen.py::test_colors_ignore_invalid | TestPyteColorsIgnoreInvalid |
| test_screen.py::test_cursor_back_last_column | TestPyteCursorBackLastColumn |
| test_screen.py::test_cursor_movement | TestPyteCursorMovement |
| test_screen.py::test_cursor_position | TestPyteCursorPosition |
| test_screen.py::test_delete_characters | TestPyteDeleteCharacters |
| test_screen.py::test_delete_lines | TestPyteDeleteLines |
| test_screen.py::test_display_complex_emoji | TestPyteDisplayComplexEmoji |
| test_screen.py::test_display_multi_char_emoji | TestPyteDisplayMultiCharEmoji |
| test_screen.py::test_display_width | TestPyteDisplayWidth |
| test_screen.py::test_draw | TestPyteDraw |
| test_screen.py::test_draw_cp437 | TestPyteDrawCP437 |
| test_screen.py::test_draw_multiple_chars | TestPyteDrawMultipleChars |
| test_screen.py::test_draw_russian | TestPyteDrawRussian |
| test_screen.py::test_draw_utf8 | TestPyteDrawUTF8 |
| test_screen.py::test_draw_width0_combining | TestPyteDrawWidth0Combining |
| test_screen.py::test_draw_width0_decawm_off | TestPyteDrawWidth0DecawmOff |
| test_screen.py::test_draw_width0_irm | TestPyteDrawWidth0IRM |
| test_screen.py::test_draw_width2 | TestPyteDrawWidth2 |
| test_screen.py::test_draw_width2_irm | TestPyteDrawWidth2IRM |
| test_screen.py::test_draw_width2_line_end | TestPyteDrawWidth2LineEnd |
| test_screen.py::test_draw_with_carriage_return | TestPyteDrawWithCarriageReturn |
| test_screen.py::test_erase | TestPyteErase |
| test_screen.py::test_erase_character | TestPyteEraseCharacter |
| test_screen.py::test_erase_in_display | TestPyteEraseInDisplay |
| test_screen.py::test_erase_in_line | TestPyteEraseInLine |
| test_screen.py::test_hide_cursor | TestPyteHideCursor |
| test_screen.py::test_index | TestPyteIndex |
| test_screen.py::test_insert_characters | TestPyteInsertCharacters |
| test_screen.py::test_insert_delete_lines_and_chars | TestPyteInsertDeleteLinesAndChars |
| test_screen.py::test_insert_lines | TestPyteInsertLines |
| test_screen.py::test_linefeed | TestPyteLinefeed |
| test_screen.py::test_linefeed_margins | TestPyteLinefeedMargins |
| test_screen.py::test_multi_attribs | TestMultiAttribs |
| test_screen.py::test_private_sgr_ignored | TestPytePrivateSGRIgnored |
| test_screen.py::test_remove_non_existent_attribute | TestPyteRemoveNonExistentAttribute |
| test_screen.py::test_report_device_attributes | TestPyteReportDeviceAttributes |
| test_screen.py::test_report_device_status | TestPyteReportDeviceStatus |
| test_screen.py::test_reset_between_attributes | TestPyteResetBetweenAttributes |
| test_screen.py::test_reset_resets_colors | TestResetResetsColors |
| test_screen.py::test_resize | TestResize |
| test_screen.py::test_resize_same | TestPyteResizeSame |
| test_screen.py::test_restore_cursor_none_saved | TestPyteRestoreCursorNoneSaved |
| test_screen.py::test_restore_cursor_out_of_bounds | TestPyteRestoreCursorOutOfBounds |
| test_screen.py::test_reverse_index | TestPyteReverseIndex |
| test_screen.py::test_save_restore_cursor | TestPyteSaveRestoreCursor |
| test_screen.py::test_screen_set_icon_name_title | TestPyteScreenSetIconNameTitle |
| test_screen.py::test_set_margins | TestPyteSetMargins |
| test_screen.py::test_set_margins_zero | TestPyteSetMarginsZero |
| test_screen.py::test_set_mode | TestPyteSetMode |
| test_screen.py::test_tab_stops | TestPyteTabStops |
| test_screen.py::test_unicode | TestPyteUnicode |
| test_stream.py::test_basic_sequences | TestBasicSequences |
| test_stream.py::test_byte_stream_define_charset | TestByteStreamDefineCharset |
| test_stream.py::test_byte_stream_define_charset_unknown | TestByteStreamDefineCharsetUnknown |
| test_stream.py::test_byte_stream_feed | TestByteStreamFeed |
| test_stream.py::test_byte_stream_select_other_charset | TestByteStreamSelectOtherCharset |
| test_stream.py::test_debug_stream | TestDebugStream |
| test_stream.py::test_dollar_skip | TestDollarSkip |
| test_stream.py::test_handler_exception | TestHandlerException |
