/* ==========================================================================
   Umbra Suara — Calculator (Music & B2B)
   Extracted from templates/division.html
   ========================================================================== */
document.addEventListener('DOMContentLoaded', function() {
(function() {
  var container = document.querySelector('.section-calculator');
  if (!container) return;
  var WA_PHONE = container.dataset.waPhone || '';
  // Tab switching
  var pathBtns = document.querySelectorAll('.calc-path-btn');
  var musicPanel = document.getElementById('pathway-music');
  var b2bPanel = document.getElementById('pathway-b2b');

  function switchPathway(name) {
    pathBtns.forEach(function(b) { b.classList.toggle('active', b.dataset.pathway === name); });
    musicPanel.classList.toggle('active', name === 'music');
    b2bPanel.classList.toggle('active', name === 'b2b');
  }

  pathBtns.forEach(function(btn) {
    btn.addEventListener('click', function() {
      switchPathway(btn.dataset.pathway);
    });
  });

  // ========== MUSIC PATHWAY (Shift-based) ==========
  var tierRadios = document.querySelectorAll('input[name="tier"]');
  var shiftsSlider = document.getElementById('shifts-slider');
  var shiftsInput = document.getElementById('shifts-input');
  var shiftsDisplay = document.getElementById('shifts-display');
  var songsSlider = document.getElementById('songs-slider');
  var songsDisplay = document.getElementById('songs-display');
  var extraRevisionsSlider = document.getElementById('extra-revisions-slider');
  var extraRevisionsDisplay = document.getElementById('extra-revisions-display');
  var expressCbMusic = document.getElementById('express-cb-music');

  var musicPreprojectLine = document.getElementById('music-preproject-line');
  var musicRecordingLine = document.getElementById('music-recording-line');
  var musicMixLine = document.getElementById('music-mix-line');
  var musicExtraRevisionsRow = document.getElementById('music-extra-revisions-row');
  var musicExtraRevisionsLine = document.getElementById('music-extra-revisions-line');
  var musicExpressRow = document.getElementById('music-express-row');
  var musicExpressLine = document.getElementById('music-express-line');
  var musicTotal = document.getElementById('music-total');

  var PRE_PROJECT_FEE = parseInt(container.dataset.preProjectFee) || 150000;
  var SHIFT_HOURS = 8;
  var MIX_HOURS_PER_SONG = 4;
  var BASE_SHIFT_COUNT = 5;
  var EXTRA_SHIFT_HOURLY = parseInt(container.dataset.extraShiftHourly) || 200000;
  var EXTRA_REVISION_RATE = parseInt(container.dataset.extraRevisionRate) || 100000;
  var EXPRESS_FEE = parseInt(container.dataset.expressFee) || 500000;

  function formatRp(amount) {
    return 'Rp ' + amount.toLocaleString('id-ID');
  }

  function updateMusic() {
    // Get tier rate
    var hourlyRate = 70000;
    var tierLabel = 'Student';
    tierRadios.forEach(function(r) {
      if (r.checked) {
        hourlyRate = parseInt(r.value);
        tierLabel = r.parentElement.querySelector('.calc-tier-label').textContent;
      }
    });

    // Effective shifts: use manual input if > 5, else slider
    var sliderVal = parseInt(shiftsSlider.value);
    var manualVal = parseInt(shiftsInput.value) || 0;
    var shifts, extraHours;
    if (manualVal > BASE_SHIFT_COUNT) {
      shifts = manualVal;
      shiftsSlider.value = BASE_SHIFT_COUNT;
      extraHours = (manualVal - BASE_SHIFT_COUNT) * SHIFT_HOURS;
    } else {
      shifts = sliderVal;
      shiftsInput.value = sliderVal;
      extraHours = 0;
    }

    shiftsDisplay.textContent = shifts + ' Shift (' + (shifts * SHIFT_HOURS) + ' Jam)';

    var songCount = parseInt(songsSlider.value);
    var songLabel = songCount + ' Lagu';
    if (songCount === 1) songLabel = '1 Lagu (Single)';
    else if (songCount >= 2 && songCount <= 7) songLabel = songCount + ' Lagu (EP)';
    else if (songCount >= 8) songLabel = songCount + ' Lagu (Album)';
    songsDisplay.textContent = songLabel;

    var extraRevisions = parseInt(extraRevisionsSlider.value);
    extraRevisionsDisplay.textContent = extraRevisions;

    var wantExpress = expressCbMusic.checked;

    // Recording: base shifts at tier rate + extra hours at penalty rate
    var paidBaseShifts = Math.min(shifts, BASE_SHIFT_COUNT);
    var baseRecordingHours = paidBaseShifts * SHIFT_HOURS;
    var recordingBase = baseRecordingHours * hourlyRate;
    var extraShiftCost = extraHours * EXTRA_SHIFT_HOURLY;
    var recordingTotal = recordingBase + extraShiftCost;

    // Mix & Master: 4 hours/song at tier rate
    var mixTotal = songCount * MIX_HOURS_PER_SONG * hourlyRate;

    // Extra revisions — only when there are songs
    if (songCount === 0) {
      extraRevisionsSlider.value = 0;
      extraRevisionsSlider.disabled = true;
      extraRevisions = 0;
      extraRevisionsDisplay.textContent = '0';
    } else {
      extraRevisionsSlider.disabled = false;
    }
    var revisionsTotal = extraRevisions * EXTRA_REVISION_RATE * songCount;

    // Express — per song, only available when there are songs
    if (songCount === 0) {
      expressCbMusic.checked = false;
      expressCbMusic.disabled = true;
    } else {
      expressCbMusic.disabled = false;
    }
    var expressTotal = wantExpress && songCount > 0 ? songCount * EXPRESS_FEE : 0;

    // Pre-project fee: only applies when there's actual work
    var hasWork = shifts > 0 || songCount > 0;
    var effectivePreProject = hasWork ? PRE_PROJECT_FEE : 0;

    // Grand total
    var grandTotal = effectivePreProject + recordingTotal + mixTotal + revisionsTotal + expressTotal;

    // Update breakdown
    var musicPreprojectRow = document.getElementById('music-preproject-row');
    if (effectivePreProject > 0) {
      musicPreprojectRow.style.display = 'flex';
    } else {
      musicPreprojectRow.style.display = 'none';
    }
    musicPreprojectLine.textContent = formatRp(effectivePreProject);
    var musicRecordingRow = document.getElementById('music-recording-row');
    var recordingLabel = document.getElementById('music-recording-label');
    if (shifts > 0) {
      musicRecordingRow.style.display = 'flex';
      recordingLabel.textContent = 'Rekaman Studio (' + shifts + ' Shift / ' + (shifts * SHIFT_HOURS) + ' Jam)';
    } else {
      musicRecordingRow.style.display = 'none';
    }
    musicRecordingLine.textContent = formatRp(recordingTotal);
    var musicMixRow = document.getElementById('music-mix-row');
    if (songCount > 0) {
      musicMixRow.style.display = 'flex';
    } else {
      musicMixRow.style.display = 'none';
    }
    musicMixLine.textContent = formatRp(mixTotal);

    var musicExtraRevisionsRow = document.getElementById('music-extra-revisions-row');
    if (songCount > 0) {
      musicExtraRevisionsRow.style.display = 'flex';
    } else {
      musicExtraRevisionsRow.style.display = 'none';
    }
    musicExtraRevisionsLine.textContent = formatRp(revisionsTotal);

    if (wantExpress) {
      musicExpressRow.style.display = 'flex';
      musicExpressLine.textContent = formatRp(expressTotal);
    } else {
      musicExpressRow.style.display = 'none';
    }

    musicTotal.textContent = formatRp(grandTotal);

    // WhatsApp pre-filled message
    var musicMsgLines = [
      'Halo Umbra Suara, saya ingin memesan produksi musik:',
      '',
      hasWork ? '• Biaya Pre-Project: ' + formatRp(effectivePreProject) : null,
      '• Paket: ' + tierLabel + ' (' + formatRp(hourlyRate) + '/jam)',
      '• Durasi Rekaman: ' + shifts + ' Shift (' + (shifts * SHIFT_HOURS) + ' Jam)',
      extraHours > 0 ? '• Kelebihan Shift: ' + extraHours + ' Jam (' + formatRp(extraShiftCost) + ')' : null,
      '• Jumlah Lagu Mix & Master: ' + songCount + ' Lagu',
      extraRevisions > 0 ? '• Revisi Tambahan: ' + extraRevisions + 'x x ' + songCount + ' lagu = ' + formatRp(revisionsTotal) : null,
      wantExpress ? '• Express 48 Jam: ' + songCount + ' lagu x ' + formatRp(EXPRESS_FEE) + ' = ' + formatRp(expressTotal) : null,
      '',
      'Total Estimasi: ' + formatRp(grandTotal)
    ];
    var msgStr = musicMsgLines.filter(function(l) { return l !== null; }).join('\n');
    var musicWA = 'https://wa.me/' + WA_PHONE + '?text=' + encodeURIComponent(msgStr);
    document.getElementById('music-cta').href = musicWA;
  }

  // Event listeners
  tierRadios.forEach(function(r) { r.addEventListener('change', updateMusic); });
  shiftsSlider.addEventListener('input', function() {
    shiftsInput.value = this.value;
    updateMusic();
  });
  shiftsInput.addEventListener('input', function() {
    var val = parseInt(this.value);
    if (isNaN(val) || val < 0) { val = 0; this.value = 0; }
    if (val <= BASE_SHIFT_COUNT) {
      shiftsSlider.value = val;
    } else {
      shiftsSlider.value = BASE_SHIFT_COUNT;
    }
    updateMusic();
  });
  songsSlider.addEventListener('input', updateMusic);
  extraRevisionsSlider.addEventListener('input', updateMusic);
  expressCbMusic.addEventListener('change', updateMusic);

  // Init
  updateMusic();

  // ========== B2B PATHWAY ==========
  var b2bProjectType = document.getElementById('b2b-project-type');
  var b2bLicensing = document.getElementById('b2b-licensing');
  var voCb = document.getElementById('vo-cb');
  var fasttrackCb = document.getElementById('fasttrack-cb');
  var b2bRevisionsSlider = document.getElementById('b2b-revisions-slider');
  var b2bRevisionsDisplay = document.getElementById('b2b-revisions-display');

  var b2bBaseLine = document.getElementById('b2b-base-line');
  var b2bSurchargeRow = document.getElementById('b2b-surcharge-row');
  var b2bSurchargeLine = document.getElementById('b2b-surcharge-line');
  var b2bVoRow = document.getElementById('b2b-vo-row');
  var b2bVoLine = document.getElementById('b2b-vo-line');
  var b2bRevisionsRow = document.getElementById('b2b-revisions-row');
  var b2bRevisionsLine = document.getElementById('b2b-revisions-line');
  var b2bMultiplierLine = document.getElementById('b2b-multiplier-line');
  var b2bTotal = document.getElementById('b2b-total');

  var VO_FLAT_FEE = parseInt(container.dataset.voFlatFee) || 1500000;
  var B2B_REVISION_RATE = parseInt(container.dataset.b2bRevisionRate) || 350000;
  var B2B_REVISION_RATE_FAST = parseInt(container.dataset.b2bRevisionRateFast) || 450000;

  function updateB2B() {
    var baseProjectPrice = parseInt(b2bProjectType.value);
    var licensingMultiplier = parseFloat(b2bLicensing.value);
    var wantVoiceOver = voCb.checked;
    var wantFastTrack = fasttrackCb.checked;
    var b2bRevisions = parseInt(b2bRevisionsSlider.value);

    b2bRevisionsDisplay.textContent = b2bRevisions;

    var FASTTRACK_FLAT = parseInt(container.dataset.fasttrackFlat) || 250000;
    var surcharge = 0;
    if (wantFastTrack) {
      surcharge = baseProjectPrice * 0.30 + FASTTRACK_FLAT;
    }

    var productionSubtotal = baseProjectPrice + surcharge;

    if (wantVoiceOver) {
      productionSubtotal += VO_FLAT_FEE;
    }

    var revisionRate = wantFastTrack ? B2B_REVISION_RATE_FAST : B2B_REVISION_RATE;
    var revisionsTotal = b2bRevisions * revisionRate;

    var grandTotal = Math.round(productionSubtotal * licensingMultiplier) + revisionsTotal;

    b2bBaseLine.textContent = 'Rp ' + baseProjectPrice.toLocaleString('id-ID');

    if (wantFastTrack) {
      b2bSurchargeRow.style.display = 'flex';
      b2bSurchargeLine.textContent = 'Rp ' + Math.round(surcharge).toLocaleString('id-ID');
    } else {
      b2bSurchargeRow.style.display = 'none';
    }

    if (wantVoiceOver) {
      b2bVoRow.style.display = 'flex';
      b2bVoLine.textContent = 'Rp ' + VO_FLAT_FEE.toLocaleString('id-ID');
    } else {
      b2bVoRow.style.display = 'none';
    }

    if (b2bRevisions > 0) {
      b2bRevisionsRow.style.display = 'flex';
      document.getElementById('b2b-revisions-label').textContent = 'Biaya Revisi (@Rp ' + Number(revisionRate).toLocaleString('id-ID') + ')';
      b2bRevisionsLine.textContent = 'Rp ' + revisionsTotal.toLocaleString('id-ID');
    } else {
      b2bRevisionsRow.style.display = 'none';
    }

    b2bMultiplierLine.textContent = licensingMultiplier.toFixed(1) + 'x';
    b2bTotal.textContent = 'Rp ' + grandTotal.toLocaleString('id-ID');

    // WhatsApp pre-filled message
    var projectLabel = b2bProjectType.options[b2bProjectType.selectedIndex].text;
    var licenseLabel = b2bLicensing.options[b2bLicensing.selectedIndex].text;
    var b2bMsgLines = [
      'Halo Umbra, saya tertarik dengan layanan audio komersial:',
      '',
      '• Tipe Proyek: ' + projectLabel,
      '• Lisensi: ' + licenseLabel,
      '• Voice-Over: ' + (wantVoiceOver ? 'Ya' : 'Tidak'),
      wantFastTrack ? '• Fast Track 72 Jam: Rp ' + Math.round(surcharge).toLocaleString('id-ID') + ' (30% + Rp250rb)' : null,
      b2bRevisions > 0 ? '• Revisi: ' + b2bRevisions + 'x @ Rp ' + Number(revisionRate).toLocaleString('id-ID') + ' = Rp ' + revisionsTotal.toLocaleString('id-ID') : null,
      '',
      'Total Estimasi: Rp ' + grandTotal.toLocaleString('id-ID')
    ];
    var msgStr = b2bMsgLines.filter(function(l) { return l !== null; }).join('\n');
    var b2bWA = 'https://wa.me/' + WA_PHONE + '?text=' + encodeURIComponent(msgStr);
    document.getElementById('b2b-cta').href = b2bWA;
  }

  b2bProjectType.addEventListener('change', updateB2B);
  b2bLicensing.addEventListener('change', updateB2B);
  voCb.addEventListener('change', updateB2B);
  fasttrackCb.addEventListener('change', updateB2B);
  b2bRevisionsSlider.addEventListener('input', updateB2B);

  // Init both
  updateMusic();
  updateB2B();
})();
});
